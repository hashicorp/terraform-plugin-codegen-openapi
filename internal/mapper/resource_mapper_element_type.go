package mapper

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func mapSchemaToElementType(schema *base.Schema) (*ir.ElementType, error) {
	oasType, err := retrieveType(schema)
	if err != nil {
		return nil, err
	}

	switch oasType {
	case oas_type_string:
		return &ir.ElementType{
			String: &ir.StringElement{},
		}, nil

	case oas_type_integer:
		return &ir.ElementType{
			Int64: &ir.Int64Element{},
		}, nil

	case oas_type_number:
		if schema.Format == oas_format_double || schema.Format == oas_format_float {
			return &ir.ElementType{
				Float64: &ir.Float64Element{},
			}, nil
		}

		return &ir.ElementType{
			Number: &ir.NumberElement{},
		}, nil

	case oas_type_boolean:
		return &ir.ElementType{
			Bool: &ir.BoolElement{},
		}, nil

	case oas_type_array:
		return mapSchemaToArrayElementType(schema)

	case oas_type_object:
		elemTypes, err := mapSchemaToObjectElementTypes(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object elem type schema proxy - %w", err)
		}

		return &ir.ElementType{
			Object: *elemTypes,
		}, nil

	default:
		return nil, fmt.Errorf("invalid schema type '%s'", oasType)
	}
}

func mapSchemaToObjectElementTypes(schema *base.Schema) (*[]ir.ObjectElement, error) {
	objectElemTypes := []ir.ObjectElement{}

	// Guarantee the order of processing
	propertyNames := sortedKeys(schema.Properties)
	for _, pName := range propertyNames {
		pProxy := schema.Properties[pName]

		pSchema, err := buildSchema(pProxy)
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		elemType, err := mapSchemaToElementType(pSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to build object property '%s' schema proxy - %w", pName, err)
		}
		objectElemTypes = append(objectElemTypes, ir.ObjectElement{
			Name:        pName,
			ElementType: elemType,
		})
	}

	return &objectElemTypes, nil
}

func mapSchemaToArrayElementType(schema *base.Schema) (*ir.ElementType, error) {
	if !schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for nested elem array, doesn't have a schema")
	}
	iSchema, err := buildSchema(schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("failed to build nested array items schema")
	}

	elementType, err := mapSchemaToElementType(iSchema)
	if err != nil {
		return nil, err
	}

	return &ir.ElementType{
		List: &ir.ListElement{
			ElementType: elementType,
		},
	}, nil
}
