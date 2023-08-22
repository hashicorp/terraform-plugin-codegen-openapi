// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildBoolResource(name string, computability schema.ComputedOptionalRequired) (mapper_resource.MapperAttribute, error) {
	result := &mapper_resource.MapperBoolAttribute{
		Name: name,
		BoolAttribute: resource.BoolAttribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	if s.Schema.Default != nil {
		staticDefault, ok := s.Schema.Default.(bool)

		if ok {
			if computability == schema.Required {
				result.ComputedOptionalRequired = schema.ComputedOptional
			}

			result.Default = &schema.BoolDefault{
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
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}

	return result, nil
}

func (s *OASSchema) BuildBoolProvider(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	return &provider.Attribute{
		Name: name,
		Bool: &provider.BoolAttribute{
			OptionalRequired:   optionalOrRequired,
			DeprecationMessage: s.GetDeprecationMessage(),
			Description:        s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildBoolElementType() (schema.ElementType, error) {
	return schema.ElementType{
		Bool: &schema.BoolType{},
	}, nil
}
