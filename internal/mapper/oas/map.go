// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func (s *OASSchema) BuildMapResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	mapSchemaProxy, ok := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	if !ok {
		return nil, fmt.Errorf("invalid map schema, expected type *base.SchemaProxy, got: %T", s.Schema.AdditionalProperties)
	}

	mapSchema, err := BuildSchema(mapSchemaProxy, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("error building map schema proxy - %w", err)
	}

	result := &resource.Attribute{
		Name: name,
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildResourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}
		result.MapNested = &resource.MapNestedAttribute{
			NestedObject: resource.NestedAttributeObject{
				Attributes: *mapAttributes,
			},
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.MapNested.Validators = s.GetMapValidators()
		}

		return result, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	result.Map = &resource.MapAttribute{
		ElementType:              elemType,
		ComputedOptionalRequired: computability,
		Description:              s.GetDescription(),
	}

	if computability != schema.Computed {
		result.Map.Validators = s.GetMapValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildMapDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	mapSchemaProxy, ok := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	if !ok {
		return nil, fmt.Errorf("invalid map schema, expected type *base.SchemaProxy, got: %T", s.Schema.AdditionalProperties)
	}

	mapSchema, err := BuildSchema(mapSchemaProxy, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("error building map schema proxy - %w", err)
	}

	result := &datasource.Attribute{
		Name: name,
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}
		result.MapNested = &datasource.MapNestedAttribute{
			NestedObject: datasource.NestedAttributeObject{
				Attributes: *mapAttributes,
			},
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.MapNested.Validators = s.GetMapValidators()
		}

		return result, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	result.Map = &datasource.MapAttribute{
		ElementType:              elemType,
		ComputedOptionalRequired: computability,
		Description:              s.GetDescription(),
	}

	if computability != schema.Computed {
		result.Map.Validators = s.GetMapValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildMapProvider(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	mapSchemaProxy, ok := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	if !ok {
		return nil, fmt.Errorf("invalid map schema, expected type *base.SchemaProxy, got: %T", s.Schema.AdditionalProperties)
	}

	mapSchema, err := BuildSchema(mapSchemaProxy, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("error building map schema proxy - %w", err)
	}

	result := &provider.Attribute{
		Name: name,
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildProviderAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}
		result.MapNested = &provider.MapNestedAttribute{
			NestedObject: provider.NestedAttributeObject{
				Attributes: *mapAttributes,
			},
			OptionalRequired: optionalOrRequired,
			Description:      s.GetDescription(),
		}

		result.MapNested.Validators = s.GetMapValidators()

		return result, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	result.Map = &provider.MapAttribute{
		ElementType:      elemType,
		OptionalRequired: optionalOrRequired,
		Description:      s.GetDescription(),
	}

	result.Map.Validators = s.GetMapValidators()

	return result, nil
}

func (s *OASSchema) BuildMapElementType() (schema.ElementType, error) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
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

func (s *OASSchema) GetMapValidators() []schema.MapValidator {
	var result []schema.MapValidator

	minProperties := s.Schema.MinProperties
	maxProperties := s.Schema.MaxProperties

	if minProperties != nil && maxProperties != nil {
		result = append(result, schema.MapValidator{
			Custom: frameworkvalidators.MapValidatorSizeBetween(*minProperties, *maxProperties),
		})
	} else if minProperties != nil {
		result = append(result, schema.MapValidator{
			Custom: frameworkvalidators.MapValidatorSizeAtLeast(*minProperties),
		})
	} else if maxProperties != nil {
		result = append(result, schema.MapValidator{
			Custom: frameworkvalidators.MapValidatorSizeAtMost(*maxProperties),
		})
	}

	return result
}
