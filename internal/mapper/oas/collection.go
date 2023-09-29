// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildCollectionResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, *PropertyError) {
	if !s.Schema.Items.IsA() {
		return nil, s.NewPropertyError(errors.New("invalid array items property, doesn't have a schema"), name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, s.NewPropertyError(fmt.Errorf("failed to build array items schema: %w", err), name)
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, propErr := itemSchema.BuildResourceAttributes()
		if propErr != nil {
			return nil, s.NestPropertyError(propErr, name)
		}

		if s.Schema.Format == util.TF_format_set {
			result := &attrmapper.ResourceSetNestedAttribute{
				Name: name,
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: objectAttributes,
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: computability,
					DeprecationMessage:       s.GetDeprecationMessage(),
					Description:              s.GetDescription(),
				},
			}

			if computability != schema.Computed {
				result.Validators = s.GetSetValidators()
			}

			return result, nil
		}

		result := &attrmapper.ResourceListNestedAttribute{
			Name: name,
			NestedObject: attrmapper.ResourceNestedAttributeObject{
				Attributes: objectAttributes,
			},
			ListNestedAttribute: resource.ListNestedAttribute{
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			},
		}

		if computability != schema.Computed {
			result.Validators = s.GetListValidators()
		}

		return result, nil
	}

	elemType, propErr := itemSchema.BuildElementType()
	if propErr != nil {
		return nil, s.NestPropertyError(propErr, name)
	}

	if s.Schema.Format == util.TF_format_set {
		result := &attrmapper.ResourceSetAttribute{
			Name: name,
			SetAttribute: resource.SetAttribute{
				ElementType:              elemType,
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			},
		}

		if computability != schema.Computed {
			result.Validators = s.GetSetValidators()
		}

		return result, nil
	}

	result := &attrmapper.ResourceListAttribute{
		Name: name,
		ListAttribute: resource.ListAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Validators = s.GetListValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildCollectionDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, *PropertyError) {
	if !s.Schema.Items.IsA() {
		return nil, s.NewPropertyError(errors.New("invalid array items property, doesn't have a schema"), name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, s.NewPropertyError(fmt.Errorf("failed to build array items schema: %w", err), name)
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, propErr := itemSchema.BuildDataSourceAttributes()
		if propErr != nil {
			return nil, s.NestPropertyError(propErr, name)
		}

		if s.Schema.Format == util.TF_format_set {

			result := &attrmapper.DataSourceSetNestedAttribute{
				Name: name,
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: objectAttributes,
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: computability,
					DeprecationMessage:       s.GetDeprecationMessage(),
					Description:              s.GetDescription(),
				},
			}

			if computability != schema.Computed {
				result.Validators = s.GetSetValidators()
			}

			return result, nil
		}

		result := &attrmapper.DataSourceListNestedAttribute{
			Name: name,
			NestedObject: attrmapper.DataSourceNestedAttributeObject{
				Attributes: objectAttributes,
			},
			ListNestedAttribute: datasource.ListNestedAttribute{
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			},
		}

		if computability != schema.Computed {
			result.Validators = s.GetListValidators()
		}

		return result, nil
	}

	elemType, propErr := itemSchema.BuildElementType()
	if propErr != nil {
		return nil, s.NestPropertyError(propErr, name)
	}

	if s.Schema.Format == util.TF_format_set {

		result := &attrmapper.DataSourceSetAttribute{
			Name: name,
			SetAttribute: datasource.SetAttribute{
				ElementType:              elemType,
				ComputedOptionalRequired: computability,
				DeprecationMessage:       s.GetDeprecationMessage(),
				Description:              s.GetDescription(),
			},
		}

		if computability != schema.Computed {
			result.Validators = s.GetSetValidators()
		}

		return result, nil
	}

	result := &attrmapper.DataSourceListAttribute{
		Name: name,
		ListAttribute: datasource.ListAttribute{
			ElementType:              elemType,
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Validators = s.GetListValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildCollectionProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, *PropertyError) {
	if !s.Schema.Items.IsA() {
		return nil, s.NewPropertyError(errors.New("invalid array items property, doesn't have a schema"), name)
	}

	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return nil, s.NewPropertyError(fmt.Errorf("failed to build array items schema: %w", err), name)
	}

	// If the items schema is a map (i.e. additionalProperties set to a schema), it cannot be a NestedAttribute
	if itemSchema.Type == util.OAS_type_object && !itemSchema.IsMap() {
		objectAttributes, propErr := itemSchema.BuildProviderAttributes()
		if propErr != nil {
			return nil, s.NestPropertyError(propErr, name)
		}

		if s.Schema.Format == util.TF_format_set {

			result := &attrmapper.ProviderSetNestedAttribute{
				Name: name,
				NestedObject: attrmapper.ProviderNestedAttributeObject{
					Attributes: objectAttributes,
				},
				SetNestedAttribute: provider.SetNestedAttribute{
					OptionalRequired:   optionalOrRequired,
					DeprecationMessage: s.GetDeprecationMessage(),
					Description:        s.GetDescription(),
					Validators:         s.GetSetValidators(),
				},
			}

			return result, nil
		}

		result := &attrmapper.ProviderListNestedAttribute{
			Name: name,
			NestedObject: attrmapper.ProviderNestedAttributeObject{
				Attributes: objectAttributes,
			},
			ListNestedAttribute: provider.ListNestedAttribute{
				OptionalRequired:   optionalOrRequired,
				DeprecationMessage: s.GetDeprecationMessage(),
				Description:        s.GetDescription(),
				Validators:         s.GetListValidators(),
			},
		}

		return result, nil
	}

	elemType, propErr := itemSchema.BuildElementType()
	if propErr != nil {
		return nil, s.NestPropertyError(propErr, name)
	}

	if s.Schema.Format == util.TF_format_set {
		result := &attrmapper.ProviderSetAttribute{
			Name: name,
			SetAttribute: provider.SetAttribute{
				ElementType:        elemType,
				OptionalRequired:   optionalOrRequired,
				DeprecationMessage: s.GetDeprecationMessage(),
				Description:        s.GetDescription(),
				Validators:         s.GetSetValidators(),
			},
		}

		return result, nil
	}

	result := &attrmapper.ProviderListAttribute{
		Name: name,
		ListAttribute: provider.ListAttribute{
			ElementType:        elemType,
			OptionalRequired:   optionalOrRequired,
			DeprecationMessage: s.GetDeprecationMessage(),
			Description:        s.GetDescription(),
			Validators:         s.GetListValidators(),
		},
	}

	return result, nil
}

func (s *OASSchema) BuildCollectionElementType() (schema.ElementType, *PropertyError) {
	if !s.Schema.Items.IsA() {
		return schema.ElementType{}, EmptyPropertyError(errors.New("invalid array type for nested elem array, doesn't have a schema"))
	}
	itemSchema, err := BuildSchema(s.Schema.Items.A, SchemaOpts{}, s.GlobalSchemaOpts)
	if err != nil {
		return schema.ElementType{}, EmptyPropertyError(errors.New("failed to build nested array items schema"))
	}

	elemType, propErr := itemSchema.BuildElementType()
	if propErr != nil {
		return schema.ElementType{}, propErr
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
