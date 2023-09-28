// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper

import (
	"errors"
	"log/slog"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
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
			rLogger.Warn("skipping resource schema mapping", "err", err)
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
	createRequestSchema, err := oas.BuildSchemaFromRequest(explorerResource.CreateOp, oas.SchemaOpts{}, oas.GlobalSchemaOpts{})
	if err != nil {
		return nil, err
	}
	createRequestAttributes, err := createRequestSchema.BuildResourceAttributes()
	if err != nil {
		return nil, err
	}

	// *********************
	// Create Response Body (optional)
	// *********************
	logger.Debug("searching for create operation response body")
	createResponseAttributes := attrmapper.ResourceAttributes{}
	createResponseSchema, err := oas.BuildSchemaFromResponse(explorerResource.CreateOp, oas.SchemaOpts{}, oas.GlobalSchemaOpts{OverrideComputability: schema.Computed})
	if err != nil {
		if errors.Is(err, oas.ErrSchemaNotFound) {
			// Demote log to INFO if there was no schema found
			logger.Info("skipping mapping of create operation response body", "err", err)
		} else {
			logger.Warn("skipping mapping of create operation response body", "err", err)
		}
	} else {
		createResponseAttributes, err = createResponseSchema.BuildResourceAttributes()
		if err != nil {
			logger.Warn("skipping mapping of create operation response body", "err", err)
		}
	}

	// *******************
	// READ Response Body (optional)
	// *******************
	logger.Debug("searching for read operation response body")
	readResponseAttributes := attrmapper.ResourceAttributes{}
	readResponseSchema, err := oas.BuildSchemaFromResponse(explorerResource.ReadOp, oas.SchemaOpts{}, oas.GlobalSchemaOpts{OverrideComputability: schema.Computed})
	if err != nil {
		if errors.Is(err, oas.ErrSchemaNotFound) {
			// Demote log to INFO if there was no schema found
			logger.Info("skipping mapping of read operation response body", "err", err)
		} else {
			logger.Warn("skipping mapping of read operation response body", "err", err)
		}
	} else {
		readResponseAttributes, err = readResponseSchema.BuildResourceAttributes()
		if err != nil {
			logger.Warn("skipping mapping of read operation response body", "err", err)
		}
	}

	// ****************
	// READ Parameters (optional)
	// ****************
	readParameterAttributes := attrmapper.ResourceAttributes{}
	if explorerResource.ReadOp != nil && explorerResource.ReadOp.Parameters != nil {
		for _, param := range explorerResource.ReadOp.Parameters {
			if param.In != util.OAS_param_path && param.In != util.OAS_param_query {
				continue
			}

			pLogger := logger.With("param", param.Name)
			schemaOpts := oas.SchemaOpts{OverrideDescription: param.Description}
			globalSchemaOpts := oas.GlobalSchemaOpts{OverrideComputability: schema.ComputedOptional}

			s, err := oas.BuildSchema(param.Schema, schemaOpts, globalSchemaOpts)
			if err != nil {
				pLogger.Warn("skipping mapping of read operation parameter", "err", err)
				continue
			}

			// Check for any aliases and replace the paramater name if found
			paramName := param.Name
			if aliasedName, ok := explorerResource.SchemaOptions.AttributeOptions.Aliases[param.Name]; ok {
				pLogger = pLogger.With("param_alias", aliasedName)
				paramName = aliasedName
			}

			parameterAttribute, err := s.BuildResourceAttribute(paramName, schema.ComputedOptional)
			if err != nil {
				pLogger.Warn("skipping mapping of read operation parameter", "err", err)
				continue
			}

			readParameterAttributes = append(readParameterAttributes, parameterAttribute)
		}
	}

	// TODO: currently, no errors can be returned from merging, but in the future we should consider raising errors/warnings for unexpected scenarios, like type mismatches between attribute schemas
	resourceAttributes, _ := createRequestAttributes.Merge(createResponseAttributes, readResponseAttributes, readParameterAttributes)

	// TODO: handle error for overrides
	resourceAttributes, _ = resourceAttributes.ApplyOverrides(explorerResource.SchemaOptions.AttributeOptions.Overrides)

	resourceSchema.Attributes = resourceAttributes.ToSpec()
	return resourceSchema, nil
}
