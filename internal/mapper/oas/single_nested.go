package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildSingleNestedResource(name string, behavior schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	objectAttributes, err := s.BuildResourceAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &resource.Attribute{
		Name: name,
		SingleNested: &resource.SingleNestedAttribute{
			Attributes:               *objectAttributes,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildSingleNestedDataSource(name string, behavior schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	objectAttributes, err := s.BuildDataSourceAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &datasource.Attribute{
		Name: name,
		SingleNested: &datasource.SingleNestedAttribute{
			Attributes:               *objectAttributes,
			ComputedOptionalRequired: behavior,
			Description:              s.GetDescription(),
		},
	}, nil
}
