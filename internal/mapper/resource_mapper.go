// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

var _ ResourceMapper = resourceMapper{}

type ResourceMapper interface {
	MapToIR() ([]resource.Resource, error)
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

func (m resourceMapper) MapToIR() ([]resource.Resource, error) {
	resourceSchemas := []resource.Resource{}

	// Guarantee the order of processing
	resourceNames := util.SortedKeys(m.resources)
	for _, name := range resourceNames {
		explorerResource := m.resources[name]

		schema, err := generateResourceSchema(explorerResource)
		if err != nil {
			log.Printf("[WARN] skipping '%s' resource schema: %s\n", name, err)
			continue
		}

		resourceSchemas = append(resourceSchemas, resource.Resource{
			Name:   name,
			Schema: schema,
		})
	}

	return resourceSchemas, nil
}

func generateResourceSchema(explorerResource explorer.Resource) (*resource.Schema, error) {
	resourceSchema := &resource.Schema{
		Attributes: []resource.Attribute{},
	}

	// ********************
	// Create Request Body (required)
	// ********************
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
	createResponseAttributes := &[]resource.Attribute{}
	createResponseSchema, err := oas.BuildSchemaFromResponse(explorerResource.CreateOp, oas.SchemaOpts{}, oas.GlobalSchemaOpts{OverrideComputability: schema.Computed})
	if err != nil && !errors.Is(err, oas.ErrSchemaNotFound) {
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
	readResponseAttributes := &[]resource.Attribute{}
	readResponseSchema, err := oas.BuildSchemaFromResponse(explorerResource.ReadOp, oas.SchemaOpts{}, oas.GlobalSchemaOpts{OverrideComputability: schema.Computed})
	if err != nil && !errors.Is(err, oas.ErrSchemaNotFound) {
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
	// TODO: Expand support for "header" and "cookie"?
	// TODO: support style + explode?
	//	- https://spec.openapis.org/oas/latest.html#style-values
	// 	- https://spec.openapis.org/oas/latest.html#style-examples
	readParameterAttributes := []resource.Attribute{}
	if explorerResource.ReadOp != nil && explorerResource.ReadOp.Parameters != nil {
		for _, param := range explorerResource.ReadOp.Parameters {
			if param.In != util.OAS_param_path && param.In != util.OAS_param_query {
				continue
			}

			schemaOpts := oas.SchemaOpts{
				OverrideDescription: param.Description,
			}

			s, err := oas.BuildSchema(param.Schema, schemaOpts, oas.GlobalSchemaOpts{OverrideComputability: schema.ComputedOptional})
			if err != nil {
				return nil, fmt.Errorf("failed to build param schema for '%s'", param.Name)
			}

			// Check for any paramater name matches and replace the name if found
			paramName := param.Name
			if matchedName, ok := explorerResource.ParameterMatches[param.Name]; ok {
				paramName = matchedName
			}

			parameterAttribute, err := s.BuildResourceAttribute(paramName, schema.ComputedOptional)
			if err != nil {
				log.Printf("[WARN] error mapping param attribute %s - %s", param.Name, err.Error())
			}

			readParameterAttributes = append(readParameterAttributes, *parameterAttribute)
		}
	}

	resourceAttributes := mergeResourceAttributes(
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
func mergeResourceAttributes(mainSlice []resource.Attribute, mergeAttributeSlices ...[]resource.Attribute) *[]resource.Attribute {
	for _, attributeSlice := range mergeAttributeSlices {

		for _, compareAttribute := range attributeSlice {
			isNewAttribute := true

			for mainIndex, mainAttribute := range mainSlice {
				if mainAttribute.Name == compareAttribute.Name {
					// Handle types that require nested merging
					if mainAttribute.SingleNested != nil && compareAttribute.SingleNested != nil {
						mergedAttributes := mergeResourceAttributes(mainAttribute.SingleNested.Attributes, compareAttribute.SingleNested.Attributes)
						mainSlice[mainIndex].SingleNested.Attributes = *mergedAttributes
					} else if mainAttribute.ListNested != nil && compareAttribute.ListNested != nil {
						mergedAttributes := mergeResourceAttributes(mainAttribute.ListNested.NestedObject.Attributes, compareAttribute.ListNested.NestedObject.Attributes)
						mainSlice[mainIndex].ListNested.NestedObject.Attributes = *mergedAttributes
					} else if mainAttribute.List != nil && compareAttribute.List != nil {
						mergedElementType := mergeElementType(mainAttribute.List.ElementType, compareAttribute.List.ElementType)
						mainSlice[mainIndex].List.ElementType = mergedElementType
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
