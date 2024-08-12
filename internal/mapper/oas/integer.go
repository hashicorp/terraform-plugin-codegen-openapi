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

func (s *OASSchema) BuildIntegerResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, *SchemaError) {
	result := &attrmapper.ResourceInt64Attribute{
		Name: name,
		Int64Attribute: resource.Int64Attribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if s.Schema.Default != nil {
		var staticDefault int64
		if err := s.Schema.Default.Decode(&staticDefault); err == nil {
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

func (s *OASSchema) BuildIntegerDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, *SchemaError) {
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

func (s *OASSchema) BuildIntegerProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, *SchemaError) {
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

func (s *OASSchema) BuildIntegerElementType() (schema.ElementType, *SchemaError) {
	return schema.ElementType{
		Int64: &schema.Int64Type{},
	}, nil
}

func (s *OASSchema) GetIntegerValidators() []schema.Int64Validator {
	var result []schema.Int64Validator

	if len(s.Schema.Enum) > 0 {
		var enum []int64

		for _, valueNode := range s.Schema.Enum {
			var value int64
			if err := valueNode.Decode(&value); err != nil {
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
