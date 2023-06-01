package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildBoolResource(name string, behavior schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	return &resource.Attribute{
		Name: name,
		Bool: &resource.BoolAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildBoolDataSource(name string, behavior schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	return &datasource.Attribute{
		Name: name,
		Bool: &datasource.BoolAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildBoolElementType() (schema.ElementType, error) {
	return schema.ElementType{
		Bool: &schema.BoolType{},
	}, nil
}
