// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildListResource(name string, behavior schema.ComputedOptionalRequired) (*resource.Attribute, error) {
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

		return &resource.Attribute{
			Name: name,
			ListNested: &resource.ListNestedAttribute{
				NestedObject: resource.NestedAttributeObject{
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

	return &resource.Attribute{
		Name: name,
		List: &resource.ListAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildListDataSource(name string, behavior schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
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

		return &datasource.Attribute{
			Name: name,
			ListNested: &datasource.ListNestedAttribute{
				NestedObject: datasource.NestedAttributeObject{
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

	return &datasource.Attribute{
		Name: name,
		List: &datasource.ListAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildListElementType() (schema.ElementType, error) {
	if !s.Schema.Items.IsA() {
		return schema.ElementType{}, fmt.Errorf("invalid array type for nested elem array, doesn't have a schema")
	}
	itemSchema, err := BuildSchema(s.Schema.Items.A)
	if err != nil {
		return schema.ElementType{}, fmt.Errorf("failed to build nested array items schema")
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return schema.ElementType{}, err
	}

	return schema.ElementType{
		List: &schema.ListType{
			ElementType: elemType,
		},
	}, nil
}
