// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildStringResource(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	return &resource.Attribute{
		Name: name,
		String: &resource.StringAttribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}, nil
}

func (s *OASSchema) BuildStringDataSource(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	return &datasource.Attribute{
		Name: name,
		String: &datasource.StringAttribute{
			ComputedOptionalRequired: computability,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}, nil
}

func (s *OASSchema) BuildStringElementType() (schema.ElementType, error) {
	return schema.ElementType{
		String: &schema.StringType{},
	}, nil
}
