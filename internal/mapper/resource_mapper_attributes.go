package mapper

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func mapSchemaToAttribute(name string, checkBehavior behaviorChecker, proxy *base.SchemaProxy) (*ir.ResourceAttribute, error) {
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema proxy - %w", err)
	}

	// Type can have multiple values or no values in an OAS, need to handle that - (╯°□°)╯︵ ┻━┻
	oasType, err := retrieveType(schema.Type)
	if err != nil {
		return nil, err
	}

	switch oasType {
	case oas_type_string:
		isSensitive := schema.Format == oas_format_password

		return &ir.ResourceAttribute{
			Name: name,
			Type: ir.ResourceAttributeType{
				String: &ir.ResourceString{
					ComputedOptionalRequired: checkBehavior(name),
					Description:              &schema.Description,
					Sensitive:                &isSensitive,
				},
			},
		}, nil

	case oas_type_integer:
		if schema.Format == oas_format_int64 {
			return &ir.ResourceAttribute{
				Name: name,
				Type: ir.ResourceAttributeType{
					Int64: &ir.ResourceInt64{
						ComputedOptionalRequired: checkBehavior(name),
						Description:              &schema.Description,
					},
				},
			}, nil
		}

		return &ir.ResourceAttribute{
			Name: name,
			Type: ir.ResourceAttributeType{
				Number: &ir.ResourceNumber{
					ComputedOptionalRequired: checkBehavior(name),
					Description:              &schema.Description,
				},
			},
		}, nil

	case oas_type_number:
		if schema.Format == oas_format_double {
			return &ir.ResourceAttribute{
				Name: name,
				Type: ir.ResourceAttributeType{
					Float64: &ir.ResourceFloat64{
						ComputedOptionalRequired: checkBehavior(name),
						Description:              &schema.Description,
					},
				},
			}, nil
		}
		return &ir.ResourceAttribute{
			Name: name,
			Type: ir.ResourceAttributeType{
				Number: &ir.ResourceNumber{
					ComputedOptionalRequired: checkBehavior(name),
					Description:              &schema.Description,
				},
			},
		}, nil

	case oas_type_boolean:
		return &ir.ResourceAttribute{
			Name: name,
			Type: ir.ResourceAttributeType{
				Bool: &ir.ResourceBool{
					ComputedOptionalRequired: checkBehavior(name),
					Description:              &schema.Description,
				},
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
			Type: ir.ResourceAttributeType{
				SingleNested: &ir.SingleNestedAttribute{
					Attributes:               *objectAttributes,
					ComputedOptionalRequired: checkBehavior(name),
					Description:              &schema.Description,
				},
			},
		}, nil

	default:
		return nil, fmt.Errorf("invalid schema type '%s'", oasType)
	}
}

func mapSchemaToObjectAttributes(proxy *base.SchemaProxy) (*[]ir.ResourceAttribute, error) {
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build object schema proxy - %w", err)
	}

	objectAttributes := []ir.ResourceAttribute{}

	// Required properties are defined a level higher then the property itself as an array ¯\_(ツ)_/¯
	behaviorChecker := propBehaviorChecker(schema.Required)

	// Guarantee the order of processing
	propertyNames := sortedKeys(schema.Properties)
	for _, name := range propertyNames {
		pProxy := schema.Properties[name]

		attribute, err := mapSchemaToAttribute(name, behaviorChecker, pProxy)
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
	arraySchema, err := schema.Items.A.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}
	oasType, err := retrieveType(arraySchema.Type)
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
			Type: ir.ResourceAttributeType{
				ListNested: &ir.ListNestedAttribute{
					NestedObject: ir.NestedObjectClass{
						Attributes: *objectAttributes,
					},
					ComputedOptionalRequired: checkBehavior(name),
					Description:              &schema.Description,
				},
			},
		}, nil
	}

	elemType, err := mapSchemaToElementType(arraySchema)
	if err != nil {
		return nil, fmt.Errorf("failed to map array elem type - %w", err)
	}

	return &ir.ResourceAttribute{
		Name: name,
		Type: ir.ResourceAttributeType{
			List: &ir.ResourceList{
				ElementType:              *elemType,
				ComputedOptionalRequired: checkBehavior(name),
				Description:              &schema.Description,
			},
		},
	}, nil
}
