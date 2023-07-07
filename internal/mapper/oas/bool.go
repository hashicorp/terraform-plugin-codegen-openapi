// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildBoolResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	result := &resource.Attribute{
		Name: name,
		Bool: &resource.BoolAttribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}

	if s.Schema.Default != nil {
		staticDefault, ok := s.Schema.Default.(bool)

		if ok {
			result.Bool.Default = &schema.BoolDefault{
				Static: &staticDefault,
			}
		}
	}

	return result, nil
}

func (s *OASSchema) BuildBoolDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	result := &datasource.Attribute{
		Name: name,
		Bool: &datasource.BoolAttribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
		},
	}

	return result, nil
}

func (s *OASSchema) BuildBoolProvider(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	return &provider.Attribute{
		Name: name,
		Bool: &provider.BoolAttribute{
			OptionalRequired: optionalOrRequired,
			Description:      s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildBoolElementType() (schema.ElementType, error) {
	return schema.ElementType{
		Bool: &schema.BoolType{},
	}, nil
}
