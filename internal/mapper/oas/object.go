package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildObjectElementType() (schema.ElementType, error) {
	objectElemTypes := []schema.ObjectAttributeType{}

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {
		pProxy := s.Schema.Properties[name]

		pSchema, err := BuildSchema(pProxy)
		if err != nil {
			return schema.ElementType{}, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		elemType, err := pSchema.BuildElementType()
		if err != nil {
			return schema.ElementType{}, fmt.Errorf("failed to create object property '%s' schema proxy - %w", name, err)
		}

		objectElemTypes = append(objectElemTypes, util.ConvertToAttributeType(name, elemType))
	}

	return schema.ElementType{
		Object: objectElemTypes,
	}, nil
}
