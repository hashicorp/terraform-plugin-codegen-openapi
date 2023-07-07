// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildNumberResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	result := &resource.Attribute{
		Name: name,
	}

	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		result.Float64 = &resource.Float64Attribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.Float64.Validators = s.GetFloatValidators()
		}

		return result, nil
	}

	result.Number = &resource.NumberAttribute{
		ComputedOptionalRequired: computability,
		Description:              s.GetDescription(),
	}

	return result, nil
}

func (s *OASSchema) BuildNumberDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	result := &datasource.Attribute{
		Name: name,
	}

	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		result.Float64 = &datasource.Float64Attribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		}

		if computability != schema.Computed {
			result.Float64.Validators = s.GetFloatValidators()
		}

		return result, nil
	}

	result.Number = &datasource.NumberAttribute{
		ComputedOptionalRequired: computability,
		Description:              s.GetDescription(),
	}

	return result, nil
}

func (s *OASSchema) BuildNumberProvider(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	result := &provider.Attribute{
		Name: name,
	}

	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		result.Float64 = &provider.Float64Attribute{
			OptionalRequired: optionalOrRequired,
			Description:      s.GetDescription(),
		}

		result.Float64.Validators = s.GetFloatValidators()

		return result, nil
	}

	result.Number = &provider.NumberAttribute{
		OptionalRequired: optionalOrRequired,
		Description:      s.GetDescription(),
	}

	return result, nil
}

func (s *OASSchema) BuildNumberElementType() (schema.ElementType, error) {
	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		return schema.ElementType{
			Float64: &schema.Float64Type{},
		}, nil
	}

	return schema.ElementType{
		Number: &schema.NumberType{},
	}, nil
}

func (s *OASSchema) GetFloatValidators() []schema.Float64Validator {
	var result []schema.Float64Validator

	if len(s.Schema.Enum) > 0 {
		var enum []float64

		for _, valueIface := range s.Schema.Enum {
			value, ok := valueIface.(float64)

			if !ok {
				// could consider error/panic here to notify developers
				continue
			}

			enum = append(enum, value)
		}

		customValidator := frameworkvalidators.Float64ValidatorOneOf(enum)

		if customValidator != nil {
			result = append(result, schema.Float64Validator{
				Custom: customValidator,
			})
		}
	}

	return result
}
