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
)

func (s *OASSchema) BuildCollectionResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	result := &resource.Attribute{
		Name: name,
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, err := itemSchema.BuildResourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		if s.Schema.Format == util.TF_format_set {
			result.SetNested = &resource.SetNestedAttribute{
				NestedObject: resource.NestedAttributeObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			}

			if computability != schema.Computed {
				result.SetNested.Validators = s.GetSetValidators()
			}

			return result, nil
		}

		result.ListNested = &resource.ListNestedAttribute{
			NestedObject: resource.NestedAttributeObject{
				Attributes: *objectAttributes,
			},
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.ListNested.Validators = s.GetListValidators()
		}

		return result, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection elem type - %w", err)
	}

	if s.Schema.Format == util.TF_format_set {
		result.Set = &resource.SetAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.Set.Validators = s.GetSetValidators()
		}

		return result, nil
	}

	result.List = &resource.ListAttribute{
		ElementType:              elemType,
		ComputedOptionalRequired: computability,
		Description:              s.GetDescription(),
	}

	if computability != schema.Computed {
		result.List.Validators = s.GetListValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildCollectionDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	result := &datasource.Attribute{
		Name: name,
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, err := itemSchema.BuildDataSourceAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		if s.Schema.Format == util.TF_format_set {
			result.SetNested = &datasource.SetNestedAttribute{
				NestedObject: datasource.NestedAttributeObject{
					Attributes: *objectAttributes,
				},
				ComputedOptionalRequired: computability,
				Description:              s.GetDescription(),
			}

			if computability != schema.Computed {
				result.SetNested.Validators = s.GetSetValidators()
			}

			return result, nil
		}

		result.ListNested = &datasource.ListNestedAttribute{
			NestedObject: datasource.NestedAttributeObject{
				Attributes: *objectAttributes,
			},
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.ListNested.Validators = s.GetListValidators()
		}

		return result, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection elem type - %w", err)
	}

	if s.Schema.Format == util.TF_format_set {
		result.Set = &datasource.SetAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.Set.Validators = s.GetSetValidators()
		}

		return result, nil
	}

	result.List = &datasource.ListAttribute{
		ElementType:              elemType,
		ComputedOptionalRequired: computability,
		Description:              s.GetDescription(),
	}

	if computability != schema.Computed {
		result.List.Validators = s.GetListValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildCollectionProvider(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	if !s.Schema.Items.IsA() {
		return nil, fmt.Errorf("invalid array type for '%s', doesn't have a schema", name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to build array items schema for '%s'", name)
	}

	result := &provider.Attribute{
		Name: name,
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, err := itemSchema.BuildProviderAttributes()
		if err != nil {
			return nil, fmt.Errorf("failed to map nested object schema proxy - %w", err)
		}

		if s.Schema.Format == util.TF_format_set {
			result.SetNested = &provider.SetNestedAttribute{
				NestedObject: provider.NestedAttributeObject{
					Attributes: *objectAttributes,
				},
				OptionalRequired: optionalOrRequired,
				Description:      s.GetDescription(),
			}

			result.SetNested.Validators = s.GetSetValidators()

			return result, nil
		}

		result.ListNested = &provider.ListNestedAttribute{
			NestedObject: provider.NestedAttributeObject{
				Attributes: *objectAttributes,
			},
			OptionalRequired: optionalOrRequired,
			Description:      s.GetDescription(),
		}

		result.ListNested.Validators = s.GetListValidators()

		return result, nil
	}

	elemType, err := itemSchema.BuildElementType()
	if err != nil {
		return nil, fmt.Errorf("failed to create collection elem type - %w", err)
	}

	if s.Schema.Format == util.TF_format_set {
		result.Set = &provider.SetAttribute{
			ElementType:      elemType,
			OptionalRequired: optionalOrRequired,
			Description:      s.GetDescription(),
		}

		result.Set.Validators = s.GetSetValidators()

		return result, nil
	}

	result.List = &provider.ListAttribute{
		ElementType:      elemType,
		OptionalRequired: optionalOrRequired,
		Description:      s.GetDescription(),
	}

	result.List.Validators = s.GetListValidators()

	return result, nil
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

func (s *OASSchema) GetListValidators() []schema.ListValidator {
	var result []schema.ListValidator

	minItems := s.Schema.MinItems
	maxItems := s.Schema.MaxItems

	if minItems != nil && maxItems != nil {
		result = append(result, schema.ListValidator{
			Custom: frameworkvalidators.ListValidatorSizeBetween(*minItems, *maxItems),
		})
	} else if minItems != nil {
		result = append(result, schema.ListValidator{
			Custom: frameworkvalidators.ListValidatorSizeAtLeast(*minItems),
		})
	} else if maxItems != nil {
		result = append(result, schema.ListValidator{
			Custom: frameworkvalidators.ListValidatorSizeAtMost(*maxItems),
		})
	}

	if s.Schema.UniqueItems != nil && *s.Schema.UniqueItems {
		result = append(result, schema.ListValidator{
			Custom: frameworkvalidators.ListValidatorUniqueValues(),
		})
	}

	return result
}

func (s *OASSchema) GetSetValidators() []schema.SetValidator {
	var result []schema.SetValidator

	minItems := s.Schema.MinItems
	maxItems := s.Schema.MaxItems

	if minItems != nil && maxItems != nil {
		result = append(result, schema.SetValidator{
			Custom: frameworkvalidators.SetValidatorSizeBetween(*minItems, *maxItems),
		})
	} else if minItems != nil {
		result = append(result, schema.SetValidator{
			Custom: frameworkvalidators.SetValidatorSizeAtLeast(*minItems),
		})
	} else if maxItems != nil {
		result = append(result, schema.SetValidator{
			Custom: frameworkvalidators.SetValidatorSizeAtMost(*maxItems),
		})
	}

	return result
}
