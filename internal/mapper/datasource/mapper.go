package datasource

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
)

var _ DataSourceMapper = dataSourceMapper{}

type DataSourceMapper interface {
	MapToIR() ([]ir.DataSource, error)
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

func (m dataSourceMapper) MapToIR() ([]ir.DataSource, error) {
	dataSourceSchemas := []ir.DataSource{}

	// Guarantee the order of processing
	dataSourceNames := util.SortedKeys(m.dataSources)
	for _, name := range dataSourceNames {
		dataSource := m.dataSources[name]

		schema, err := generateDataSourceSchema(dataSource)
		if err != nil {
			log.Printf("[WARN] skipping '%s' data source schema: %s\n", name, err)
			continue
		}

		dataSourceSchemas = append(dataSourceSchemas, ir.DataSource{
			Name:   name,
			Schema: *schema,
		})
	}

	return dataSourceSchemas, nil
}

func generateDataSourceSchema(dataSource explorer.DataSource) (*ir.DataSourceSchema, error) {
	dataSourceSchema := &ir.DataSourceSchema{
		Attributes: []ir.DataSourceAttribute{},
	}

	// ****************
	// READ Parameters (optional)
	// ****************
	readParameterAttributes := []ir.DataSourceAttribute{}
	if dataSource.ReadOp != nil && dataSource.ReadOp.Parameters != nil {
		for _, param := range dataSource.ReadOp.Parameters {
			// TODO: Filter specific "in" values? (query, path, cookies (lol)) - https://spec.openapis.org/oas/latest.html#fixed-fields-9
			s, err := schema.BuildSchema(param.Schema)
			if err != nil {
				return nil, fmt.Errorf("failed to build param schema for '%s'", param.Name)
			}

			behavior := ir.ComputedOptional
			if param.Required {
				behavior = ir.Required
			}

			// TODO: schema description is preferred over param.Description. This should probably be changed
			parameterAttribute, err := s.BuildDataSourceAttribute(param.Name, behavior)
			if err != nil {
				log.Printf("[WARN] error mapping param attribute %s - %s", param.Name, err.Error())
			}

			readParameterAttributes = append(readParameterAttributes, *parameterAttribute)
		}
	}

	// ********************
	// READ Response Body (required)
	// ********************
	readResponseSchema, err := schema.BuildSchemaFromResponse(dataSource.ReadOp)
	if err != nil {
		return nil, err
	}
	readResponseAttributes, err := readResponseSchema.BuildDataSourceAttributes()
	if err != nil {
		return nil, err
	}

	dataSourceAttributes := deepMergeAttributes(
		readParameterAttributes,
		*readResponseAttributes,
	)

	dataSourceSchema.Attributes = *dataSourceAttributes
	return dataSourceSchema, nil
}

// mainSlice takes priority in the merge, will have each subsequent mergeAttributeSlice applied in sequence
// - No re-ordering of the mainSlice is done, so will append new attributes as they are encountered
func deepMergeAttributes(mainSlice []ir.DataSourceAttribute, mergeAttributeSlices ...[]ir.DataSourceAttribute) *[]ir.DataSourceAttribute {
	for _, attributeSlice := range mergeAttributeSlices {

		for _, compareAttribute := range attributeSlice {
			isNewAttribute := true

			for mainIndex, mainAttribute := range mainSlice {
				if mainAttribute.Name == compareAttribute.Name {
					// Handle types that require nested merging
					if mainAttribute.SingleNested != nil && compareAttribute.SingleNested != nil {
						mergedAttributes := deepMergeAttributes(mainAttribute.SingleNested.Attributes, compareAttribute.SingleNested.Attributes)
						mainSlice[mainIndex].SingleNested.Attributes = *mergedAttributes
					} else if mainAttribute.ListNested != nil && compareAttribute.ListNested != nil {
						mergedAttributes := deepMergeAttributes(mainAttribute.ListNested.NestedObject.Attributes, compareAttribute.ListNested.NestedObject.Attributes)
						mainSlice[mainIndex].ListNested.NestedObject.Attributes = *mergedAttributes
					} else if mainAttribute.List != nil && compareAttribute.List != nil {
						mergedElementType := deepMergeElementType(&mainAttribute.List.ElementType, &compareAttribute.List.ElementType)
						mainSlice[mainIndex].List.ElementType = *mergedElementType
					}

					isNewAttribute = false
					break
				}
			}

			if isNewAttribute {
				// Add this back to the original slice to avoid adding duplicate attributes from different mergeAttributeSlices
				mainSlice = append(mainSlice, compareAttribute)
			}
		}

	}
	return &mainSlice
}

func deepMergeElementType(mainElementType *ir.ElementType, mergeElementType *ir.ElementType) *ir.ElementType {
	if mainElementType.List != nil && mergeElementType.List != nil {
		mainElementType.List.ElementType = deepMergeElementType(mainElementType.List.ElementType, mergeElementType.List.ElementType)
	} else if mainElementType.Object != nil && mergeElementType.Object != nil {
		objectElemTypes := deepMergeObjectElementTypes(mainElementType.Object, mergeElementType.Object)
		mainElementType.Object = objectElemTypes
	}

	return mainElementType
}

func deepMergeObjectElementTypes(mainObject []ir.ObjectElement, mergeObject []ir.ObjectElement) []ir.ObjectElement {
	for _, compareElemType := range mergeObject {
		isNewElemType := true

		for mainIndex, mainElemType := range mainObject {
			if mainElemType.Name == compareElemType.Name {
				mergedElementType := deepMergeElementType(mainElemType.ElementType, compareElemType.ElementType)
				mainObject[mainIndex].ElementType = mergedElementType

				isNewElemType = false
				break
			}
		}

		if isNewElemType {
			mainObject = append(mainObject, compareElemType)
		}
	}

	return mainObject
}
