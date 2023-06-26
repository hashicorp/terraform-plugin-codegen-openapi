// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func (s *OASSchema) BuildMapResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	mapSchemaProxy, ok := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	if !ok {
		return nil, fmt.Errorf("invalid map schema, expected type *base.SchemaProxy, got: %T", s.Schema.AdditionalProperties)
	}

	mapSchema, err := BuildSchema(mapSchemaProxy, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("error building map schema proxy - %w", err)
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildResourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}
		return &resource.Attribute{
			Name: name,
			MapNested: &resource.MapNestedAttribute{
				NestedObject: resource.NestedAttributeObject{
					Attributes: *mapAttributes,
				},
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	return &resource.Attribute{
		Name: name,
		Map: &resource.MapAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildMapDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	mapSchemaProxy, ok := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	if !ok {
		return nil, fmt.Errorf("invalid map schema, expected type *base.SchemaProxy, got: %T", s.Schema.AdditionalProperties)
	}

	mapSchema, err := BuildSchema(mapSchemaProxy, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("error building map schema proxy - %w", err)
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}
		return &datasource.Attribute{
			Name: name,
			MapNested: &datasource.MapNestedAttribute{
				NestedObject: datasource.NestedAttributeObject{
					Attributes: *mapAttributes,
				},
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	return &datasource.Attribute{
		Name: name,
		Map: &datasource.MapAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildMapElementType() (schema.ElementType, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema `additionalProperties` can
	// also be a boolean (which we should ignore and map to an ObjectType), so calling functions should call s.IsMap() first.
	mapSchemaProxy, ok := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	if !ok {
		return schema.ElementType{}, fmt.Errorf("invalid map schema, expected type *base.SchemaProxy, got: %T", s.Schema.AdditionalProperties)
	}

	mapSchema, err := BuildSchema(mapSchemaProxy, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return schema.ElementType{}, fmt.Errorf("error building map schema proxy - %w", err)
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return schema.ElementType{}, fmt.Errorf("failed to create map elem type - %w", err)
	}

	return schema.ElementType{
		Map: &schema.MapType{
			ElementType: elemType,
		},
	}, nil
}
