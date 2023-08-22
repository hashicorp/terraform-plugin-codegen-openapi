// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildResourceAttributes() (mapper_resource.MapperAttributes, error) {
	objectAttributes := mapper_resource.MapperAttributes{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, SchemaOpts{}, s.GlobalSchemaOpts)
		if err != nil {
			return nil, err
		}

		attribute, err := pSchema.BuildResourceAttribute(name, s.GetComputability(name))
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema - %w", name, err)
		}

		objectAttributes = append(objectAttributes, attribute)
	}

	return objectAttributes, nil
}

func (s *OASSchema) BuildResourceAttribute(name string, computability schema.ComputedOptionalRequired) (mapper_resource.MapperAttribute, error) {
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
		return s.BuildCollectionResource(name, computability)
	case util.OAS_type_object:
		if s.IsMap() {
			return s.BuildMapResource(name, computability)
		}
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
		pSchema, err := BuildSchema(pProxy, SchemaOpts{}, s.GlobalSchemaOpts)
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
		return s.BuildCollectionDataSource(name, computability)
	case util.OAS_type_object:
		if s.IsMap() {
			return s.BuildMapDataSource(name, computability)
		}
		return s.BuildSingleNestedDataSource(name, computability)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *OASSchema) BuildProviderAttributes() (*[]provider.Attribute, error) {
	objectAttributes := []provider.Attribute{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, SchemaOpts{}, s.GlobalSchemaOpts)
		if err != nil {
			return nil, err
		}

		attribute, err := pSchema.BuildProviderAttribute(name, s.GetOptionalOrRequired(name))
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema - %w", name, err)
		}

		objectAttributes = append(objectAttributes, *attribute)
	}

	return &objectAttributes, nil
}

func (s *OASSchema) BuildProviderAttribute(name string, optionalOrRequired schema.OptionalRequired) (*provider.Attribute, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringProvider(name, optionalOrRequired)
	case util.OAS_type_integer:
		return s.BuildIntegerProvider(name, optionalOrRequired)
	case util.OAS_type_number:
		return s.BuildNumberProvider(name, optionalOrRequired)
	case util.OAS_type_boolean:
		return s.BuildBoolProvider(name, optionalOrRequired)
	case util.OAS_type_array:
		return s.BuildCollectionProvider(name, optionalOrRequired)
	case util.OAS_type_object:
		if s.IsMap() {
			return s.BuildMapProvider(name, optionalOrRequired)
		}
		return s.BuildSingleNestedProvider(name, optionalOrRequired)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}
