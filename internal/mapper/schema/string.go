package schema

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
)

func (s *OASSchema) BuildStringResource(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	return &ir.ResourceAttribute{
		Name: name,
		String: &ir.ResourceStringAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}, nil
}

func (s *OASSchema) BuildStringDataSource(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	return &ir.DataSourceAttribute{
		Name: name,
		String: &ir.DataSourceStringAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
			Sensitive:                s.IsSensitive(),
		},
	}, nil
}

func (s *OASSchema) BuildStringElementType() (*ir.ElementType, error) {
	return &ir.ElementType{
		String: &ir.StringElement{},
	}, nil
}
