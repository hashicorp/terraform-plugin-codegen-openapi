// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildIntegerResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	result := &resource.Attribute{
		Name: name,
		Int64: &resource.Int64Attribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Int64.Validators = s.GetIntegerValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildIntegerDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	result := &datasource.Attribute{
		Name: name,
		Int64: &datasource.Int64Attribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Int64.Validators = s.GetIntegerValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildIntegerProvider(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	result := &provider.Attribute{
		Name: name,
		Int64: &provider.Int64Attribute{
			OptionalRequired: optionalOrRequired,
			Description:      s.GetDescription(),
		},
	}

	result.Int64.Validators = s.GetIntegerValidators()

	return result, nil
}

func (s *OASSchema) BuildIntegerElementType() (schema.ElementType, error) {
	return schema.ElementType{
		Int64: &schema.Int64Type{},
	}, nil
}

func (s *OASSchema) GetIntegerValidators() []schema.Int64Validator {
	var result []schema.Int64Validator

	if len(s.Schema.Enum) > 0 {
		var enum []int64

		for _, valueIface := range s.Schema.Enum {
			value, ok := valueIface.(int64)

			if !ok {
				// could consider error/panic here to notify developers
				continue
			}

			enum = append(enum, value)
		}

		customValidator := frameworkvalidators.Int64ValidatorOneOf(enum)

		if customValidator != nil {
			result = append(result, schema.Int64Validator{
				Custom: customValidator,
			})
		}
	}

	minimum := s.Schema.Minimum
	maximum := s.Schema.Maximum

	if minimum != nil && maximum != nil {
		result = append(result, schema.Int64Validator{
			Custom: frameworkvalidators.Int64ValidatorBetween(int64(*minimum), int64(*maximum)),
		})
	} else if minimum != nil {
		result = append(result, schema.Int64Validator{
			Custom: frameworkvalidators.Int64ValidatorAtLeast(int64(*minimum)),
		})
	} else if maximum != nil {
		result = append(result, schema.Int64Validator{
			Custom: frameworkvalidators.Int64ValidatorAtMost(int64(*maximum)),
		})
	}

	return result
}
