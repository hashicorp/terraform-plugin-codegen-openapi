// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
)

func (s *OASSchema) BuildStringResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, *SchemaError) {
	result := &attrmapper.ResourceStringAttribute{
		Name: name,
		StringAttribute: resource.StringAttribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}

	if s.Schema.Default != nil {
		var staticDefault string
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
			if computability == schema.Required {
				result.ComputedOptionalRequired = schema.ComputedOptional
			}

			result.Default = &schema.StringDefault{
				Static: &staticDefault,
			}
		}
	}

	if computability != schema.Computed {
		result.Validators = s.GetStringValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildStringDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, *SchemaError) {
	result := &attrmapper.DataSourceStringAttribute{
		Name: name,
		StringAttribute: datasource.StringAttribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}

	if computability != schema.Computed {
		result.Validators = s.GetStringValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildStringProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, *SchemaError) {
	result := &attrmapper.ProviderStringAttribute{
		Name: name,
		StringAttribute: provider.StringAttribute{
			OptionalRequired:   optionalOrRequired,
			DeprecationMessage: s.GetDeprecationMessage(),
			Description:        s.GetDescription(),
			Sensitive:          s.IsSensitive(),
			Validators:         s.GetStringValidators(),
		},
	}

	return result, nil
}

func (s *OASSchema) BuildStringElementType() (schema.ElementType, *SchemaError) {
	return schema.ElementType{
		String: &schema.StringType{},
	}, nil
}

func (s *OASSchema) GetStringValidators() []schema.StringValidator {
	var result []schema.StringValidator

	if len(s.Schema.Enum) > 0 {
		var enum []string

		for _, valueNode := range s.Schema.Enum {
			var value string
			if err := valueNode.Decode(&value); err != nil {
				// could consider error/panic here to notify developers
				continue
			}

			enum = append(enum, value)
		}

		customValidator := frameworkvalidators.StringValidatorOneOf(enum)

		if customValidator != nil {
			result = append(result, schema.StringValidator{
				Custom: customValidator,
			})
		}
	}

	minLength := s.Schema.MinLength
	maxLength := s.Schema.MaxLength

	if minLength != nil && maxLength != nil {
		result = append(result, schema.StringValidator{
			Custom: frameworkvalidators.StringValidatorLengthBetween(*minLength, *maxLength),
		})
	} else if minLength != nil {
		result = append(result, schema.StringValidator{
			Custom: frameworkvalidators.StringValidatorLengthAtLeast(*minLength),
		})
	} else if maxLength != nil {
		result = append(result, schema.StringValidator{
			Custom: frameworkvalidators.StringValidatorLengthAtMost(*maxLength),
		})
	}

	if s.Schema.Pattern != "" {
		result = append(result, schema.StringValidator{
			// Friendly regex message could be added later via configuration or
			// custom annotation.
			Custom: frameworkvalidators.StringValidatorRegexMatches(s.Schema.Pattern, ""),
		})
	}

	return result
}
