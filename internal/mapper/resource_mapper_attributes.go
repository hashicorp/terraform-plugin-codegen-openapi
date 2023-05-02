package mapper

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func mapSchemaToAttribute(name string, checkBehavior behaviorChecker, proxy *base.SchemaProxy) (*ir.ResourceAttribute, error) {
	schema, err := buildSchema(proxy)
	if err != nil {
		return nil, err
	}

	oasType, err := retrieveType(schema)
	if err != nil {
		return nil, err
	}

	switch oasType {
	case oas_type_string:
		isSensitive := schema.Format == oas_format_password

		return &ir.ResourceAttribute{
			Name: name,
			String: &ir.ResourceStringAttribute{
				ComputedOptionalRequired: checkBehavior(name, schema),
				Description:              &schema.Description,
				Sensitive:                &isSensitive,
			},
		}, nil

	case oas_type_integer:
		return &ir.ResourceAttribute{
			Name: name,
			Int64: &ir.ResourceInt64Attribute{
				ComputedOptionalRequired: checkBehavior(name, schema),
				Description:              &schema.Description,
			},
		}, nil

	case oas_type_number:
		if schema.Format == oas_format_double || schema.Format == oas_format_float {
			return &ir.ResourceAttribute{
				Name: name,
				Float64: &ir.ResourceFloat64Attribute{
					ComputedOptionalRequired: checkBehavior(name, schema),
					Description:              &schema.Description,
				},
			}, nil
		}

		return &ir.ResourceAttribute{
			Name: name,
			Number: &ir.ResourceNumberAttribute{
				ComputedOptionalRequired: checkBehavior(name, schema),
				Description:              &schema.Description,
			},
		}, nil

	case oas_type_boolean:
		return &ir.ResourceAttribute{
			Name: name,
			Bool: &ir.ResourceBoolAttribute{
				ComputedOptionalRequired: checkBehavior(name, schema),
				Description:              &schema.Description,
			},
		}, nil
	case oas_type_array:
		return mapSchemaToArrayAttribute(name, checkBehavior, schema)

	case oas_type_object:
		objectAttributes, err := mapSchemaToObjectAttributes(proxy)
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		return &ir.ResourceAttribute{
			Name: name,
			SingleNested: &ir.ResourceSingleNestedAttribute{
				Attributes:               *objectAttributes,
				ComputedOptionalRequired: checkBehavior(name, schema),
				Description:              &schema.Description,
			},
		}, nil

	default:
		return nil, fmt.Errorf("invalid schema type '%s'", oasType)
	}
}

func mapSchemaToObjectAttributes(proxy *base.SchemaProxy) (*[]ir.ResourceAttribute, error) {
	schema, err := buildSchema(proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to build object schema proxy - %w", err)
	}

	objectAttributes := []ir.ResourceAttribute{}

	// Required properties are defined a level higher then the property itself as an array ¯\_(ツ)_/¯
	checkBehavior := propBehaviorChecker(schema.Required)

	// Guarantee the order of processing
	propertyNames := sortedKeys(schema.Properties)
	for _, name := range propertyNames {
		pProxy := schema.Properties[name]

		attribute, err := mapSchemaToAttribute(name, checkBehavior, pProxy)
		if err != nil {
			return nil, fmt.Errorf("failed to map object property '%s' schema - %w", name, err)
		}
		objectAttributes = append(objectAttributes, *attribute)
	}

	return &objectAttributes, nil
}

func mapSchemaToArrayAttribute(name string, checkBehavior behaviorChecker, schema *base.Schema) (*ir.ResourceAttribute, error) {
	if !schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	arraySchema, err := buildSchema(schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}
	oasType, err := retrieveType(arraySchema)
	if err != nil {
		return nil, err
	}

	if oasType == oas_type_object {
		objectAttributes, err := mapSchemaToObjectAttributes(schema.Items.A)
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		return &ir.ResourceAttribute{
			Name: name,
			ListNested: &ir.ResourceListNestedAttribute{
				NestedObject: ir.ResourceAttributeNestedObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: checkBehavior(name, schema),
				Description:              &schema.Description,
			},
		}, nil
	}

	elemType, err := mapSchemaToElementType(arraySchema)
	if err != nil {
		return nil, fmt.Errorf("failed to map array elem type - %w", err)
	}

	return &ir.ResourceAttribute{
		Name: name,
		List: &ir.ResourceListAttribute{
			ElementType:              *elemType,
			ComputedOptionalRequired: checkBehavior(name, schema),
			Description:              &schema.Description,
		},
	}, nil
}
