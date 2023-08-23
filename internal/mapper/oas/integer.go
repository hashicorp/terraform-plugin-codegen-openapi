// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildIntegerResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, error) {
	result := &attrmapper.ResourceInt64Attribute{
		Name: name,
		Int64Attribute: resource.Int64Attribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if s.Schema.Default != nil {
		staticDefault, ok := s.Schema.Default.(int64)

		if ok {
			if computability == schema.Required {
				result.ComputedOptionalRequired = schema.ComputedOptional
			}

			result.Default = &schema.Int64Default{
				Static: &staticDefault,
			}
		}
	}

	if computability != schema.Computed {
		result.Validators = s.GetIntegerValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildIntegerDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, error) {
	result := &attrmapper.DataSourceInt64Attribute{
		Name: name,
		Int64Attribute: datasource.Int64Attribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if computability != schema.Computed {
		result.Validators = s.GetIntegerValidators()
	}

	return result, nil
}

func (s *OASSchema) BuildIntegerProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, error) {
	result := &attrmapper.ProviderInt64Attribute{
		Name: name,
		Int64Attribute: provider.Int64Attribute{
			OptionalRequired:   optionalOrRequired,
			DeprecationMessage: s.GetDeprecationMessage(),
			Description:        s.GetDescription(),
			Validators:         s.GetIntegerValidators(),
		},
	}

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
