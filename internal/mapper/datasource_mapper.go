// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper

import (
	"log/slog"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

var _ DataSourceMapper = dataSourceMapper{}

type DataSourceMapper interface {
	MapToIR(*slog.Logger) ([]datasource.DataSource, error)
}

type dataSourceMapper struct {
	dataSources map[string]explorer.DataSource
	//nolint:unused // Might be useful later!
	cfg config.Config
}

func NewDataSourceMapper(dataSources map[string]explorer.DataSource, cfg config.Config) DataSourceMapper {
	return dataSourceMapper{
		dataSources: dataSources,
		cfg:         cfg,
	}
}

func (m dataSourceMapper) MapToIR(logger *slog.Logger) ([]datasource.DataSource, error) {
	dataSourceSchemas := []datasource.DataSource{}

	// Guarantee the order of processing
	dataSourceNames := util.SortedKeys(m.dataSources)
	for _, name := range dataSourceNames {
		dataSource := m.dataSources[name]
		dLogger := logger.With("data_source", name)

		schema, err := generateDataSourceSchema(dLogger, dataSource)
		if err != nil {
			dLogger.Warn("skipping data source schema mapping", "err", err)
			continue
		}

		dataSourceSchemas = append(dataSourceSchemas, datasource.DataSource{
			Name:   name,
			Schema: schema,
		})
	}

	return dataSourceSchemas, nil
}

func generateDataSourceSchema(logger *slog.Logger, dataSource explorer.DataSource) (*datasource.Schema, error) {
	dataSourceSchema := &datasource.Schema{
		Attributes: []datasource.Attribute{},
	}

	// ********************
	// READ Response Body (required)
	// ********************
	logger.Debug("searching for read operation response body")
	readResponseSchema, err := oas.BuildSchemaFromResponse(dataSource.ReadOp, oas.SchemaOpts{}, oas.GlobalSchemaOpts{OverrideComputability: schema.Computed})
	if err != nil {
		return nil, err
	}
	readResponseAttributes, err := readResponseSchema.BuildDataSourceAttributes()
	if err != nil {
		return nil, err
	}

	// ****************
	// READ Parameters (optional)
	// ****************
	// TODO: Expand support for "header" and "cookie"?
	// TODO: support style + explode?
	//	- https://spec.openapis.org/oas/latest.html#style-values
	// 	- https://spec.openapis.org/oas/latest.html#style-examples
	readParameterAttributes := attrmapper.DataSourceAttributes{}
	if dataSource.ReadOp != nil && dataSource.ReadOp.Parameters != nil {
		for _, param := range dataSource.ReadOp.Parameters {
			if param.In != util.OAS_param_path && param.In != util.OAS_param_query {
				continue
			}

			pLogger := logger.With("param", param.Name)
			schemaOpts := oas.SchemaOpts{OverrideDescription: param.Description}

			s, err := oas.BuildSchema(param.Schema, schemaOpts, oas.GlobalSchemaOpts{})
			if err != nil {
				pLogger.Warn("skipping mapping of read operation parameter", "err", err)
				continue
			}

			computability := schema.ComputedOptional
			if param.Required {
				computability = schema.Required
			}

			// Check for any aliases and replace the paramater name if found
			paramName := param.Name
			if aliasedName, ok := dataSource.SchemaOptions.AttributeOptions.Aliases[param.Name]; ok {
				pLogger = pLogger.With("param_alias", aliasedName)
				paramName = aliasedName
			}

			parameterAttribute, err := s.BuildDataSourceAttribute(paramName, computability)
			if err != nil {
				pLogger.Warn("skipping mapping of read operation parameter", "err", err)
				continue
			}

			readParameterAttributes = append(readParameterAttributes, parameterAttribute)
		}
	}

	// TODO: currently, no errors can be returned from merging, but in the future we should consider raising errors/warnings for unexpected scenarios, like type mismatches between attribute schemas
	dataSourceAttributes, _ := readParameterAttributes.Merge(readResponseAttributes)

	// TODO: handle error for overrides
	dataSourceAttributes, _ = dataSourceAttributes.ApplyOverrides(dataSource.SchemaOptions.AttributeOptions.Overrides)

	dataSourceSchema.Attributes = dataSourceAttributes.ToSpec()
	return dataSourceSchema, nil
}
