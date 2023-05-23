package schema

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
)

func (s *OASSchema) BuildObjectElementType() (*ir.ElementType, error) {
	objectElemTypes := []ir.ObjectElement{}

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {
		pProxy := s.Schema.Properties[name]

		pSchema, err := BuildSchema(pProxy)
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		elemType, err := pSchema.BuildElementType()
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema proxy - %w", name, err)
		}

		objectElemTypes = append(objectElemTypes, ir.ObjectElement{
			Name:        name,
			ElementType: elemType,
		})
	}

	return &ir.ElementType{
		Object: objectElemTypes,
	}, nil
}
