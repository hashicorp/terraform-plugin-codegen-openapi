// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"

	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	lowmodel "github.com/pb33f/libopenapi/datamodel/low"
	lowbase "github.com/pb33f/libopenapi/datamodel/low/base"
	low "github.com/pb33f/libopenapi/datamodel/low/v3"
)

var _ Explorer = configExplorer{}

type configExplorer struct {
	spec   high.Document
	config config.Config
}

// A ConfigExplorer will use an additional config file to identify resource and data source operations in a provided
// OpenAPIv3 spec. This additional config file will provide information such as:
//   - Create/Read/Update/Delete endpoints/URLs (schema will be automatically grabbed via request/response body and parameters in mapper)
//   - Resource + Data Source names
func NewConfigExplorer(spec high.Document, cfg config.Config) Explorer {
	return configExplorer{
		spec:   spec,
		config: cfg,
	}
}

func (e configExplorer) FindProvider() (Provider, error) {
	foundProvider := Provider{
		Name: e.config.Provider.Name,
	}

	if e.config.Provider.SchemaRef == "" {
		return foundProvider, nil
	}

	schemaProxy, err := extractSchemaProxy(e.spec, e.config.Provider.SchemaRef)
	if err != nil {
		return Provider{}, fmt.Errorf("error extracting provider schema from ref: %w", err)
	}
	foundProvider.SchemaProxy = schemaProxy
	foundProvider.Ignores = e.config.Provider.Ignores

	return foundProvider, nil
}

func (e configExplorer) FindResources() (map[string]Resource, error) {
	resources := map[string]Resource{}
	var errResult error

	for name, resourceConfig := range e.config.Resources {
		createOp, err := extractOp(e.spec.Paths, resourceConfig.Create)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s.create': %w", name, err))
			continue
		}
		readOp, err := extractOp(e.spec.Paths, resourceConfig.Read)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s.read': %w", name, err))
			continue
		}
		updateOp, err := extractOp(e.spec.Paths, resourceConfig.Update)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s.update': %w", name, err))
			continue
		}
		deleteOp, err := extractOp(e.spec.Paths, resourceConfig.Delete)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s.delete': %w", name, err))
			continue
		}

		commonParameters, err := extractCommonParameters(e.spec.Paths, resourceConfig.Read.Path)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s' common parameters: %w", name, err))
			continue
		}

		resources[name] = Resource{
			CreateOp:         createOp,
			ReadOp:           readOp,
			UpdateOp:         updateOp,
			DeleteOp:         deleteOp,
			CommonParameters: commonParameters,
			SchemaOptions:    extractSchemaOptions(resourceConfig.SchemaOptions),
		}
	}

	return resources, errResult
}

func (e configExplorer) FindDataSources() (map[string]DataSource, error) {
	dataSources := map[string]DataSource{}
	var errResult error

	for name, dataSourceConfig := range e.config.DataSources {
		readOp, err := extractOp(e.spec.Paths, dataSourceConfig.Read)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s.read': %w", name, err))
			continue
		}

		commonParameters, err := extractCommonParameters(e.spec.Paths, dataSourceConfig.Read.Path)
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("failed to extract '%s' common parameters: %w", name, err))
			continue
		}

		dataSources[name] = DataSource{
			ReadOp:           readOp,
			CommonParameters: commonParameters,
			SchemaOptions:    extractSchemaOptions(dataSourceConfig.SchemaOptions),
		}
	}
	return dataSources, errResult
}

func extractOp(paths *high.Paths, oasLocation *config.OpenApiSpecLocation) (*high.Operation, error) {
	// No need to search OAS if not defined
	if oasLocation == nil {
		return nil, nil
	}

	if paths == nil || paths.PathItems == nil || paths.PathItems.GetOrZero(oasLocation.Path) == nil {
		return nil, fmt.Errorf("path '%s' not found in OpenAPI spec", oasLocation.Path)
	}

	pathItem, _ := paths.PathItems.Get(oasLocation.Path)

	switch strings.ToLower(oasLocation.Method) {
	case low.PostLabel:
		return pathItem.Post, nil
	case low.GetLabel:
		return pathItem.Get, nil
	case low.PutLabel:
		return pathItem.Put, nil
	case low.DeleteLabel:
		return pathItem.Delete, nil
	case low.PatchLabel:
		return pathItem.Patch, nil
	case low.OptionsLabel:
		return pathItem.Options, nil
	case low.HeadLabel:
		return pathItem.Head, nil
	case low.TraceLabel:
		return pathItem.Trace, nil
	default:
		return nil, fmt.Errorf("method '%s' not found at OpenAPI path '%s'", oasLocation.Method, oasLocation.Path)
	}
}

func extractCommonParameters(paths *high.Paths, path string) ([]*high.Parameter, error) {
	// No need to search OAS if not defined
	if paths.PathItems.GetOrZero(path) == nil {
		return nil, fmt.Errorf("path '%s' not found in OpenAPI spec", path)
	}

	pathItem, _ := paths.PathItems.Get(path)

	return pathItem.Parameters, nil
}

func extractSchemaProxy(document high.Document, componentRef string) (*highbase.SchemaProxy, error) {
	// find the reference using the root document.Index
	indexRef := document.Index.FindComponentInRoot(componentRef)
	if indexRef == nil {
		return nil, fmt.Errorf("unable to find reference: %s", componentRef)
	}

	// build low-level schema using YAML node
	var lowSchema lowbase.Schema
	err := lowmodel.BuildModel(indexRef.Node, &lowSchema)
	if err != nil {
		return nil, fmt.Errorf("error building low-level schema: %w", err)
	}

	// populate low-level schema, using root document.Index for resolving
	err = lowSchema.Build(context.TODO(), indexRef.Node, document.Index)
	if err != nil {
		return nil, fmt.Errorf("error populating low-level schema: %w", err)
	}

	// build high-level schema from low-level schema
	highSchema := highbase.NewSchema(&lowSchema)

	// wrap in a schema proxy for mapping with `oas` package
	return highbase.CreateSchemaProxy(highSchema), nil
}

func extractSchemaOptions(cfgSchemaOpts config.SchemaOptions) SchemaOptions {
	return SchemaOptions{
		Ignores: cfgSchemaOpts.Ignores,
		AttributeOptions: AttributeOptions{
			Aliases:   cfgSchemaOpts.AttributeOptions.Aliases,
			Overrides: extractOverrides(cfgSchemaOpts.AttributeOptions.Overrides),
		},
	}
}

func extractOverrides(cfgOverrides map[string]config.Override) map[string]Override {
	overrides := make(map[string]Override, len(cfgOverrides))
	for key, cfgOverride := range cfgOverrides {
		overrides[key] = Override{
			Description:              cfgOverride.Description,
			ComputedOptionalRequired: cfgOverride.ComputedOptionalRequired,
		}
	}

	return overrides
}
