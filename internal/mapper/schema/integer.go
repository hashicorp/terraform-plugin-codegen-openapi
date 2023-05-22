package schema

import "github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"

func (s *OASSchema) BuildIntegerResource(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	return &ir.ResourceAttribute{
		Name: name,
		Int64: &ir.ResourceInt64Attribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildIntegerDataSource(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	return &ir.DataSourceAttribute{
		Name: name,
		Int64: &ir.DataSourceInt64Attribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildIntegerElementType() (*ir.ElementType, error) {
	return &ir.ElementType{
		Int64: &ir.Int64Element{},
	}, nil
}
