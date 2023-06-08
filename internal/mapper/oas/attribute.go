// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildResourceAttributes() (*[]resource.Attribute, error) {
	objectAttributes := []resource.Attribute{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, s.GlobalSchemaOpts)
		if err != nil {
			return nil, err
		}

		attribute, err := pSchema.BuildResourceAttribute(name, s.GetComputability(name))
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema - %w", name, err)
		}

		objectAttributes = append(objectAttributes, *attribute)
	}

	return &objectAttributes, nil
}

func (s *OASSchema) BuildResourceAttribute(name string, computability schema.ComputedOptionalRequired) (*resource.Attribute, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringResource(name, computability)
	case util.OAS_type_integer:
		return s.BuildIntegerResource(name, computability)
	case util.OAS_type_number:
		return s.BuildNumberResource(name, computability)
	case util.OAS_type_boolean:
		return s.BuildBoolResource(name, computability)
	case util.OAS_type_array:
		return s.BuildListResource(name, computability)
	case util.OAS_type_object:
		return s.BuildSingleNestedResource(name, computability)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *OASSchema) BuildDataSourceAttributes() (*[]datasource.Attribute, error) {
	objectAttributes := []datasource.Attribute{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, s.GlobalSchemaOpts)
		if err != nil {
			return nil, err
		}

		attribute, err := pSchema.BuildDataSourceAttribute(name, s.GetComputability(name))
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema - %w", name, err)
		}

		objectAttributes = append(objectAttributes, *attribute)
	}

	return &objectAttributes, nil
}

func (s *OASSchema) BuildDataSourceAttribute(name string, computability schema.ComputedOptionalRequired) (*datasource.Attribute, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringDataSource(name, computability)
	case util.OAS_type_integer:
		return s.BuildIntegerDataSource(name, computability)
	case util.OAS_type_number:
		return s.BuildNumberDataSource(name, computability)
	case util.OAS_type_boolean:
		return s.BuildBoolDataSource(name, computability)
	case util.OAS_type_array:
		return s.BuildListDataSource(name, computability)
	case util.OAS_type_object:
		return s.BuildSingleNestedDataSource(name, computability)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}
