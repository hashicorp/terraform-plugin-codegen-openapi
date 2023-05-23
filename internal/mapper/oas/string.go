package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildStringResource(name string, behavior schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	return &resource.Attribute{
		Name: name,
		String: &resource.StringAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}, nil
}

func (s *OASSchema) BuildStringDataSource(name string, behavior schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	return &datasource.Attribute{
		Name: name,
		String: &datasource.StringAttribute{
			ComputedOptionalRequired: behavior,
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
