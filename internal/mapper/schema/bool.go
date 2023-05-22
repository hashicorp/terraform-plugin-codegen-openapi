package schema

import "github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"

func (s *OASSchema) BuildBoolResource(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	return &ir.ResourceAttribute{
		Name: name,
		Bool: &ir.ResourceBoolAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildBoolDataSource(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	return &ir.DataSourceAttribute{
		Name: name,
		Bool: &ir.DataSourceBoolAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildBoolElementType() (*ir.ElementType, error) {
	return &ir.ElementType{
		Bool: &ir.BoolElement{},
	}, nil
}
