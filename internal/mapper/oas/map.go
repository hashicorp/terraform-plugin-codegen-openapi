// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func (s *OASSchema) BuildMapResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, error) {
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

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildResourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}
		result := &attrmapper.ResourceMapNestedAttribute{
			Name: name,
			NestedObject: attrmapper.ResourceNestedAttributeObject{
				Attributes: mapAttributes,
			},
			MapNestedAttribute: resource.MapNestedAttribute{
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			},
		}

		if computability != schema.Computed {
			result.Validators = s.GetMapValidators()
		}

		return result, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	result := &attrmapper.ResourceMapAttribute{
		Name: name,
		MapAttribute: resource.MapAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Validators = s.GetMapValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildMapDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, error) {
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

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		result := &attrmapper.DataSourceMapNestedAttribute{
			Name: name,
			NestedObject: attrmapper.DataSourceNestedAttributeObject{
				Attributes: mapAttributes,
			},
			MapNestedAttribute: datasource.MapNestedAttribute{
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			},
		}

		if computability != schema.Computed {
			result.Validators = s.GetMapValidators()
		}

		return result, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	result := &attrmapper.DataSourceMapAttribute{
		Name: name,
		MapAttribute: datasource.MapAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Validators = s.GetMapValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildMapProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, error) {
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

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildProviderAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
		}

		result := &attrmapper.ProviderMapNestedAttribute{
			Name: name,
			NestedObject: attrmapper.ProviderNestedAttributeObject{
				Attributes: mapAttributes,
			},
			MapNestedAttribute: provider.MapNestedAttribute{
				OptionalRequired:   optionalOrRequired,
				DeprecationMessage: s.GetDeprecationMessage(),
				Description:        s.GetDescription(),
				Validators:         s.GetMapValidators(),
			},
		}

		return result, nil
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create map elem type - %w", err)
	}

	result := &attrmapper.ProviderMapAttribute{
		Name: name,
		MapAttribute: provider.MapAttribute{
			ElementType:        elemType,
			OptionalRequired:   optionalOrRequired,
			DeprecationMessage: s.GetDeprecationMessage(),
			Description:        s.GetDescription(),
			Validators:         s.GetMapValidators(),
		},
	}

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
