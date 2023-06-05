// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildNumberResource(name string, behavior schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		return &resource.Attribute{
			Name: name,
			Float64: &resource.Float64Attribute{
				ComputedOptionalRequired: behavior,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	return &resource.Attribute{
		Name: name,
		Number: &resource.NumberAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildNumberDataSource(name string, behavior schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		return &datasource.Attribute{
			Name: name,
			Float64: &datasource.Float64Attribute{
				ComputedOptionalRequired: behavior,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	return &datasource.Attribute{
		Name: name,
		Number: &datasource.NumberAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
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
