package mapper

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func mapSchemaToElementType(schema *base.Schema) (*ir.ElementType, error) {
	oasType, err := retrieveType(schema.Type)
	if err != nil {
		return nil, err
	}

	switch oasType {
	case oas_type_string:
		return &ir.ElementType{
			String: &ir.ElementTypeString{},
		}, nil

	case oas_type_integer:
		if schema.Format == oas_format_int64 {
			return &ir.ElementType{
				Int64: &ir.ElementTypeInt64{},
			}, nil
		}

		return &ir.ElementType{
			Number: &ir.ElementTypeNumber{},
		}, nil

	case oas_type_number:
		if schema.Format == oas_format_double {
			return &ir.ElementType{
				Float64: &ir.ElementTypeFloat64{},
			}, nil
		}

		return &ir.ElementType{
			Number: &ir.ElementTypeNumber{},
		}, nil

	case oas_type_boolean:
		return &ir.ElementType{
			Bool: &ir.ElementTypeBool{},
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

func mapSchemaToObjectElementTypes(schema *base.Schema) (*[]ir.ElementTypeObject, error) {
	objectElemTypes := []ir.ElementTypeObject{}

	// Guarantee the order of processing
	propertyNames := sortedKeys(schema.Properties)
	for _, pName := range propertyNames {
		pProxy := schema.Properties[pName]

		pSchema, err := pProxy.BuildSchema()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		elemType, err := mapSchemaToElementType(pSchema)
		if err != nil {
			return nil, fmt.Errorf("failed to build object property '%s' schema proxy - %w", pName, err)
		}
		objectElemTypes = append(objectElemTypes, ir.ElementTypeObject{
			Name: pName,
			Type: *elemType,
		})
	}

	return &objectElemTypes, nil
}

func mapSchemaToArrayElementType(schema *base.Schema) (*ir.ElementType, error) {
	if !schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for nested elem array, doesn't have a schema")
	}
	iSchema, err := schema.Items.A.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested array items schema")
	}

	elementType, err := mapSchemaToElementType(iSchema)
	if err != nil {
		return nil, err
	}

	return &ir.ElementType{
		List: &ir.ElementTypeList{
			ElementType: *elementType,
		},
	}, nil
}
