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

func (s *OASSchema) BuildCollectionResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, err := itemSchema.BuildResourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		if s.Schema.Format == util.TF_format_set {
			return &resource.Attribute{
				Name: name,
				SetNested: &resource.SetNestedAttribute{
					NestedObject: resource.NestedAttributeObject{
						Attributes: *objectAttributes,
					},
					ComputedOptionalRequired: computability,
					Description:              s.GetDescription(),
				},
			}, nil
		}

		return &resource.Attribute{
			Name: name,
			ListNested: &resource.ListNestedAttribute{
				NestedObject: resource.NestedAttributeObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection elem type - %w", err)
	}

	if s.Schema.Format == util.TF_format_set {
		return &resource.Attribute{
			Name: name,
			Set: &resource.SetAttribute{
				ElementType:              elemType,
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	return &resource.Attribute{
		Name: name,
		List: &resource.ListAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildCollectionDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, err := itemSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		if s.Schema.Format == util.TF_format_set {
			return &datasource.Attribute{
				Name: name,
				SetNested: &datasource.SetNestedAttribute{
					NestedObject: datasource.NestedAttributeObject{
						Attributes: *objectAttributes,
					},
					ComputedOptionalRequired: computability,
					Description:              s.GetDescription(),
				},
			}, nil
		}

		return &datasource.Attribute{
			Name: name,
			ListNested: &datasource.ListNestedAttribute{
				NestedObject: datasource.NestedAttributeObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection elem type - %w", err)
	}

	if s.Schema.Format == util.TF_format_set {
		return &datasource.Attribute{
			Name: name,
			Set: &datasource.SetAttribute{
				ElementType:              elemType,
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	return &datasource.Attribute{
		Name: name,
		List: &datasource.ListAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildCollectionElementType() (schema.ElementType, error) {
	if !s.Schema.Items.IsA() {
		return schema.ElementType{}, fmt.Errorf("invalid array type for nested elem array, doesn't have a schema")
	}
	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return schema.ElementType{}, fmt.Errorf("failed to build nested array items schema")
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return schema.ElementType{}, err
	}

	if s.Schema.Format == util.TF_format_set {
		return schema.ElementType{
			Set: &schema.SetType{
				ElementType: elemType,
			},
		}, nil
	}

	return schema.ElementType{
		List: &schema.ListType{
			ElementType: elemType,
		},
	}, nil
}
