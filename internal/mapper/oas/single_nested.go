// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildSingleNestedResource(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, error) {
	objectAttributes, err := s.BuildResourceAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &attrmapper.ResourceSingleNestedAttribute{
		Name:       name,
		Attributes: objectAttributes,
		SingleNestedAttribute: resource.SingleNestedAttribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildSingleNestedDataSource(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, error) {
	objectAttributes, err := s.BuildDataSourceAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &attrmapper.DataSourceSingleNestedAttribute{
		Name:       name,
		Attributes: objectAttributes,
		SingleNestedAttribute: datasource.SingleNestedAttribute{
			ComputedOptionalRequired: computability,
			DeprecationMessage:       s.GetDeprecationMessage(),
			Description:              s.GetDescription(),
		},
	}, nil
}

func (s *OASSchema) BuildSingleNestedProvider(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, error) {
	objectAttributes, err := s.BuildProviderAttributes()
	if err != nil {
		return nil, fmt.Errorf("failed to build nested object schema proxy - %w", err)
	}

	return &attrmapper.ProviderSingleNestedAttribute{
		Name:       name,
		Attributes: objectAttributes,
		SingleNestedAttribute: provider.SingleNestedAttribute{
			OptionalRequired:   optionalOrRequired,
			DeprecationMessage: s.GetDeprecationMessage(),
			Description:        s.GetDescription(),
		},
	}, nil
}
