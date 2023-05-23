package schema

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
)

func (s *OASSchema) BuildNumberResource(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		return &ir.ResourceAttribute{
			Name: name,
			Float64: &ir.ResourceFloat64Attribute{
				ComputedOptionalRequired: behavior,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	return &ir.ResourceAttribute{
		Name: name,
		Number: &ir.ResourceNumberAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildNumberDataSource(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		return &ir.DataSourceAttribute{
			Name: name,
			Float64: &ir.DataSourceFloat64Attribute{
				ComputedOptionalRequired: behavior,
				Description:              s.GetDescription(),
			},
		}, nil
	}

	return &ir.DataSourceAttribute{
		Name: name,
		Number: &ir.DataSourceNumberAttribute{
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildNumberElementType() (*ir.ElementType, error) {
	if s.Format == util.OAS_format_double || s.Format == util.OAS_format_float {
		return &ir.ElementType{
			Float64: &ir.Float64Element{},
		}, nil
	}

	return &ir.ElementType{
		Number: &ir.NumberElement{},
	}, nil
}
