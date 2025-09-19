// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/log"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

var _ ResourceMapper = resourceMapper{}

type ResourceMapper interface {
	MapToIR(*slog.Logger) ([]resource.Resource, error)
}

type resourceMapper struct {
	resources map[string]explorer.Resource
	//nolint:unused // Might be useful later!
	cfg config.Config
}

func NewResourceMapper(resources map[string]explorer.Resource, cfg config.Config) ResourceMapper {
	return resourceMapper{
		resources: resources,
		cfg:       cfg,
	}
}

func (m resourceMapper) MapToIR(logger *slog.Logger) ([]resource.Resource, error) {
	resourceSchemas := []resource.Resource{}

	// Guarantee the order of processing
	resourceNames := util.SortedKeys(m.resources)
	for _, name := range resourceNames {
		explorerResource := m.resources[name]
		rLogger := logger.With("resource", name)

		schema, err := generateResourceSchema(rLogger, explorerResource)
		if err != nil {
			log.WarnLogOnError(rLogger, err, "skipping resource schema mapping")
			continue
		}

		resourceSchemas = append(resourceSchemas, resource.Resource{
			Name:   name,
			Schema: schema,
		})
	}

	return resourceSchemas, nil
}

func generateResourceSchema(logger *slog.Logger, explorerResource explorer.Resource) (*resource.Schema, error) {
	resourceSchema := &resource.Schema{
		Attributes: []resource.Attribute{},
	}

	// ********************
	// Create Request Body (required)
	// ********************
	logger.Debug("searching for create operation request body")

	schemaOpts := oas.SchemaOpts{
		Ignores: explorerResource.SchemaOptions.Ignores,
	}
	createRequestSchema, err := oas.BuildSchemaFromRequest(explorerResource.CreateOp, schemaOpts, oas.GlobalSchemaOpts{})
	if err != nil {
		return nil, err
	}
	createRequestAttributes, schemaErr := createRequestSchema.BuildResourceAttributes()
	if schemaErr != nil {
		return nil, schemaErr
	}

	// *********************
	// Create Response Body (optional)
	// *********************
	logger.Debug("searching for create operation response body")

	createResponseAttributes := attrmapper.ResourceAttributes{}
	schemaOpts = oas.SchemaOpts{
		Ignores: explorerResource.SchemaOptions.Ignores,
	}
	globalSchemaOpts := oas.GlobalSchemaOpts{
		OverrideComputability: schema.Computed,
	}
	createResponseSchema, err := oas.BuildSchemaFromResponse(explorerResource.CreateOp, schemaOpts, globalSchemaOpts)
	if err != nil {
		if errors.Is(err, oas.ErrSchemaNotFound) {
			// Demote log to INFO if there was no schema found
			logger.Info("skipping mapping of create operation response body", "err", err)
		} else {
			logger.Warn("skipping mapping of create operation response body", "err", err)
		}
	} else {
		createResponseAttributes, schemaErr = createResponseSchema.BuildResourceAttributes()
		if schemaErr != nil {
			log.WarnLogOnError(logger, schemaErr, "skipping mapping of create operation response body")
		}
	}

	// *******************
	// READ Response Body (optional)
	// *******************
	logger.Debug("searching for read operation response body")

	readResponseAttributes := attrmapper.ResourceAttributes{}

	schemaOpts = oas.SchemaOpts{
		Ignores: explorerResource.SchemaOptions.Ignores,
	}
	globalSchemaOpts = oas.GlobalSchemaOpts{
		OverrideComputability: schema.Computed,
	}
	readResponseSchema, err := oas.BuildSchemaFromResponse(explorerResource.ReadOp, schemaOpts, globalSchemaOpts)
	if err != nil {
		if errors.Is(err, oas.ErrSchemaNotFound) {
			// Demote log to INFO if there was no schema found
			logger.Info("skipping mapping of read operation response body", "err", err)
		} else {
			logger.Warn("skipping mapping of read operation response body", "err", err)
		}
	} else {
		readResponseAttributes, schemaErr = readResponseSchema.BuildResourceAttributes()
		if schemaErr != nil {
			log.WarnLogOnError(logger, schemaErr, "skipping mapping of read operation response body")
		}
	}

	// ****************
	// READ Parameters (optional)
	// ****************
	for _, param := range explorerResource.ReadOpParameters() {
		if param.In != util.OAS_param_path && param.In != util.OAS_param_query {
			continue
		}

		pLogger := logger.With("param", param.Name)
		schemaOpts := oas.SchemaOpts{
			Ignores:             explorerResource.SchemaOptions.Ignores,
			OverrideDescription: param.Description,
		}
		globalSchemaOpts := oas.GlobalSchemaOpts{OverrideComputability: schema.ComputedOptional}

		s, schemaErr := oas.BuildSchema(param.Schema, schemaOpts, globalSchemaOpts)
		if schemaErr != nil {
			log.WarnLogOnError(pLogger, schemaErr, "skipping mapping of read operation parameter")
			continue
		}

		// Check for any aliases and replace the paramater name if found
		paramName := param.Name
		if aliasedName, ok := explorerResource.SchemaOptions.AttributeOptions.Aliases[param.Name]; ok {
			pLogger = pLogger.With("param_alias", aliasedName)
			paramName = aliasedName
		}

		if s.IsPropertyIgnored(paramName) {
			continue
		}

		paramPath := strings.Split(paramName, ".")

		parameterAttribute, schemaErr := s.BuildResourceAttribute(paramPath[len(paramPath)-1], schema.ComputedOptional)
		if schemaErr != nil {
			log.WarnLogOnError(pLogger, schemaErr, "skipping mapping of read operation parameter")
			continue
		}

		// TODO: currently, no errors can be returned from merging, but in the future we should consider raising errors/warnings for unexpected scenarios, like type mismatches between attribute schemas
		readResponseAttributes, _ = readResponseAttributes.MergeAttribute(paramPath, parameterAttribute, schema.ComputedOptional)
	}

	// TODO: currently, no errors can be returned from merging, but in the future we should consider raising errors/warnings for unexpected scenarios, like type mismatches between attribute schemas
	resourceAttributes, _ := createRequestAttributes.Merge(createResponseAttributes, readResponseAttributes)

	// TODO: handle error for overrides
	resourceAttributes, _ = resourceAttributes.ApplyOverrides(explorerResource.SchemaOptions.AttributeOptions.Overrides)

	resourceSchema.Attributes = resourceAttributes.ToSpec()
	return resourceSchema, nil
}
