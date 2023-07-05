// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildStringResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	result := &resource.Attribute{
		Name: name,
		String: &resource.StringAttribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}

	if computability != schema.Computed {
		result.String.Validators = s.GetStringValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildStringDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	result := &datasource.Attribute{
		Name: name,
		String: &datasource.StringAttribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}

	if computability != schema.Computed {
		result.String.Validators = s.GetStringValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildStringElementType() (schema.ElementType, error) {
	return schema.ElementType{
		String: &schema.StringType{},
	}, nil
}

func (s *OASSchema) GetStringValidators() []schema.StringValidator {
	var result []schema.StringValidator

	if len(s.Schema.Enum) > 0 {
		var enum []string

		for _, valueIface := range s.Schema.Enum {
			value, ok := valueIface.(string)

			if !ok {
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
