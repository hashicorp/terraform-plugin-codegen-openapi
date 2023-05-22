package schema

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/util"
)

func (s *OASSchema) BuildListResource(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	if itemSchema.Type == util.OAS_type_object {
		objectAttributes, err := itemSchema.BuildResourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		return &ir.ResourceAttribute{
			Name: name,
			ListNested: &ir.ResourceListNestedAttribute{
				NestedObject: ir.ResourceAttributeNestedObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: behavior,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create list elem type - %w", err)
	}

	return &ir.ResourceAttribute{
		Name: name,
		List: &ir.ResourceListAttribute{
			ElementType:              *elemType,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildListDataSource(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	if itemSchema.Type == util.OAS_type_object {
		objectAttributes, err := itemSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		return &ir.DataSourceAttribute{
			Name: name,
			ListNested: &ir.DataSourceListNestedAttribute{
				NestedObject: ir.DataSourceAttributeNestedObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: behavior,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create list elem type - %w", err)
	}

	return &ir.DataSourceAttribute{
		Name: name,
		List: &ir.DataSourceListAttribute{
			ElementType:              *elemType,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildListElementType() (*ir.ElementType, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for nested elem array, doesn't have a schema")
	}
	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return nil, fmt.Errorf("failed to build nested array items schema")
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, err
	}

	return &ir.ElementType{
		List: &ir.ListElement{
			ElementType: elemType,
		},
	}, nil
}
