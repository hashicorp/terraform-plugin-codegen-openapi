// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"errors"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildMapResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, *SchemaError) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	if !s.IsMap() {
		return nil, s.SchemaErrorFromProperty(errors.New("invalid map, additionalProperties doesn't have a valid schema"), name)
	}

	schemaOpts := SchemaOpts{
		Ignores: s.SchemaOpts.Ignores,
	}
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A, schemaOpts, s.GlobalSchemaOpts)
	if err != nil {
		return nil, s.NestSchemaError(err, name)
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildResourceAttributes()
		if err != nil {
			return nil, s.NestSchemaError(err, name)
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
		return nil, s.NestSchemaError(err, name)
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

func (s *OASSchema) BuildMapDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, *SchemaError) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	if !s.IsMap() {
		return nil, s.SchemaErrorFromProperty(errors.New("invalid map, additionalProperties doesn't have a valid schema"), name)
	}

	schemaOpts := SchemaOpts{
		Ignores: s.SchemaOpts.Ignores,
	}
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A, schemaOpts, s.GlobalSchemaOpts)
	if err != nil {
		return nil, s.NestSchemaError(err, name)
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, s.NestSchemaError(err, name)
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
		return nil, s.NestSchemaError(err, name)
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

func (s *OASSchema) BuildMapProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, *SchemaError) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an SingleNestedAttribute), so calling functions should call s.IsMap() first.
	if !s.IsMap() {
		return nil, s.SchemaErrorFromProperty(errors.New("invalid map, additionalProperties doesn't have a valid schema"), name)
	}

	schemaOpts := SchemaOpts{
		Ignores: s.SchemaOpts.Ignores,
	}
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A, schemaOpts, s.GlobalSchemaOpts)
	if err != nil {
		return nil, s.NestSchemaError(err, name)
	}

	if mapSchema.Type == util.OAS_type_object {
		mapAttributes, err := mapSchema.BuildProviderAttributes()
		if err != nil {
			return nil, s.NestSchemaError(err, name)
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
		return nil, s.NestSchemaError(err, name)
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

func (s *OASSchema) BuildMapElementType() (schema.ElementType, *SchemaError) {
	// Maps are detected as `type: object`, with an `additionalProperties` field that is a schema. `additionalProperties` can
	// also be a boolean (which we should ignore and map to an ObjectType), so calling functions should call s.IsMap() first.
	if !s.IsMap() {
		return schema.ElementType{}, SchemaErrorFromNode(errors.New("invalid map, additionalProperties doesn't have a valid schema"), s.Schema, AdditionalProperties)
	}

	schemaOpts := SchemaOpts{
		Ignores: s.SchemaOpts.Ignores,
	}
	mapSchema, err := BuildSchema(s.Schema.AdditionalProperties.A, schemaOpts, s.GlobalSchemaOpts)
	if err != nil {
		return schema.ElementType{}, err
	}

	elemType, err := mapSchema.BuildElementType()
	if err != nil {
		return schema.ElementType{}, err
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
