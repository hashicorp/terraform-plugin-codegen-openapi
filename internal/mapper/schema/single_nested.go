package schema

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
)

// TODO: move to object go file?
func (s *OASSchema) BuildSingleNestedResource(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	objectAttributes, err := s.BuildResourceAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &ir.ResourceAttribute{
		Name: name,
		SingleNested: &ir.ResourceSingleNestedAttribute{
			Attributes:               *objectAttributes,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildSingleNestedDataSource(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	objectAttributes, err := s.BuildDataSourceAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &ir.DataSourceAttribute{
		Name: name,
		SingleNested: &ir.DataSourceSingleNestedAttribute{
			Attributes:               *objectAttributes,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}
