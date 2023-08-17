// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer

import (
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

	return foundProvider, nil
}

func (e configExplorer) FindResources() (map[string]Resource, error) {
	resources := map[string]Resource{}
	for name, resourceConfig := range e.config.Resources {
		// TODO: should we throw an error if an invalid or non-existent path/methods are given?
		resources[name] = Resource{
			CreateOp:      extractOp(e.spec.Paths, resourceConfig.Create),
			ReadOp:        extractOp(e.spec.Paths, resourceConfig.Read),
			UpdateOp:      extractOp(e.spec.Paths, resourceConfig.Update),
			DeleteOp:      extractOp(e.spec.Paths, resourceConfig.Delete),
			SchemaOptions: extractSchemaOptions(resourceConfig.SchemaOptions),
		}
	}

	return resources, nil
}

func (e configExplorer) FindDataSources() (map[string]DataSource, error) {
	dataSources := map[string]DataSource{}
	for name, dataSourceConfig := range e.config.DataSources {
		dataSources[name] = DataSource{
			ReadOp:        extractOp(e.spec.Paths, dataSourceConfig.Read),
			SchemaOptions: extractSchemaOptions(dataSourceConfig.SchemaOptions),
		}
	}
	return dataSources, nil
}

func extractOp(paths *high.Paths, oasLocation *config.OpenApiSpecLocation) *high.Operation {
	if oasLocation == nil || paths == nil || paths.PathItems == nil || paths.PathItems[oasLocation.Path] == nil {
		return nil
	}

	switch strings.ToLower(oasLocation.Method) {
	case low.PostLabel:
		return paths.PathItems[oasLocation.Path].Post
	case low.GetLabel:
		return paths.PathItems[oasLocation.Path].Get
	case low.PutLabel:
		return paths.PathItems[oasLocation.Path].Put
	case low.DeleteLabel:
		return paths.PathItems[oasLocation.Path].Delete
	case low.PatchLabel:
		return paths.PathItems[oasLocation.Path].Patch
	case low.OptionsLabel:
		return paths.PathItems[oasLocation.Path].Options
	case low.HeadLabel:
		return paths.PathItems[oasLocation.Path].Head
	case low.TraceLabel:
		return paths.PathItems[oasLocation.Path].Trace
	default:
		return nil
	}
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
	err = lowSchema.Build(indexRef.Node, document.Index)
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
		AttributeOptions: AttributeOptions{
			Aliases: cfgSchemaOpts.AttributeOptions.Aliases,
		},
	}
}
