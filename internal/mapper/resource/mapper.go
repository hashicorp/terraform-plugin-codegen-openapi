package resource

import (
	"errors"
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/config"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/schema"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/util"
	"log"
)

var _ ResourceMapper = resourceMapper{}

type ResourceMapper interface {
	MapToIR() ([]ir.Resource, error)
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

func (m resourceMapper) MapToIR() ([]ir.Resource, error) {
	resourceSchemas := []ir.Resource{}

	// Guarantee the order of processing
	resourceNames := util.SortedKeys(m.resources)
	for _, name := range resourceNames {
		resource := m.resources[name]

		schema, err := generateResourceSchema(resource)
		if err != nil {
			log.Printf("[WARN] skipping '%s' resource schema: %s\n", name, err)
			continue
		}

		resourceSchemas = append(resourceSchemas, ir.Resource{
			Name:   name,
			Schema: *schema,
		})
	}

	return resourceSchemas, nil
}

func generateResourceSchema(resource explorer.Resource) (*ir.ResourceSchema, error) {
	resourceSchema := &ir.ResourceSchema{
		Attributes: []ir.ResourceAttribute{},
	}

	// ********************
	// Create Request Body (required)
	// ********************
	createRequestSchema, err := schema.BuildSchemaFromRequest(resource.CreateOp)
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
	createResponseAttributes := &[]ir.ResourceAttribute{}
	createResponseSchema, err := schema.BuildSchemaFromResponse(resource.CreateOp)
	if err != nil && !errors.Is(err, schema.ErrSchemaNotFound) {
		return nil, err
	} else if createResponseSchema != nil {
		createResponseAttributes, err = createResponseSchema.BuildResourceAttributes()
		if err != nil {
			return nil, err
		}
	}

	// *******************
	// READ Response Body (optional)
	// *******************
	readResponseAttributes := &[]ir.ResourceAttribute{}
	readResponseSchema, err := schema.BuildSchemaFromResponse(resource.ReadOp)
	if err != nil && !errors.Is(err, schema.ErrSchemaNotFound) {
		return nil, err
	} else if readResponseSchema != nil {
		readResponseAttributes, err = readResponseSchema.BuildResourceAttributes()
		if err != nil {
			return nil, err
		}
	}

	// ****************
	// READ Parameters (optional)
	// ****************
	readParameterAttributes := []ir.ResourceAttribute{}
	if resource.ReadOp != nil && resource.ReadOp.Parameters != nil {
		for _, param := range resource.ReadOp.Parameters {
			// TODO: Filter specific "in" values? (query, path, cookies (lol)) - https://spec.openapis.org/oas/latest.html#fixed-fields-9
			s, err := schema.BuildSchema(param.Schema)
			if err != nil {
				return nil, fmt.Errorf("failed to build param schema for '%s'", param.Name)
			}

			// TODO: schema description is preferred over param.Description. This should probably be changed
			parameterAttribute, err := s.BuildResourceAttribute(param.Name, ir.ComputedOptional)
			if err != nil {
				log.Printf("[WARN] error mapping param attribute %s - %s", param.Name, err.Error())
			}

			readParameterAttributes = append(readParameterAttributes, *parameterAttribute)
		}
	}

	resourceAttributes := deepMergeAttributes(
		*createRequestAttributes,
		*createResponseAttributes,
		*readResponseAttributes,
		readParameterAttributes,
	)

	resourceSchema.Attributes = *resourceAttributes
	return resourceSchema, nil
}

// mainSlice takes priority in the merge, will have each subsequent mergeAttributeSlice applied in sequence
// - No re-ordering of the mainSlice is done, so will append new attributes as they are encountered
func deepMergeAttributes(mainSlice []ir.ResourceAttribute, mergeAttributeSlices ...[]ir.ResourceAttribute) *[]ir.ResourceAttribute {
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
