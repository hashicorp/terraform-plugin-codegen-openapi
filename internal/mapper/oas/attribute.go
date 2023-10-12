// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildResourceAttributes() (attrmapper.ResourceAttributes, *PropertyError) {
	objectAttributes := attrmapper.ResourceAttributes{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, SchemaOpts{}, s.GlobalSchemaOpts)
		if err != nil {
			return nil, s.NewPropertyError(err, name)
		}

		attribute, propErr := pSchema.BuildResourceAttribute(name, s.GetComputability(name))
		if propErr != nil {
			return nil, propErr
		}

		objectAttributes = append(objectAttributes, attribute)
	}

	return objectAttributes, nil
}

func (s *OASSchema) BuildResourceAttribute(name string, computability schema.ComputedOptionalRequired) (attrmapper.ResourceAttribute, *PropertyError) {
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
		return nil, s.NewPropertyError(fmt.Errorf("invalid schema type '%s'", s.Type), name)
	}
}

func (s *OASSchema) BuildDataSourceAttributes() (attrmapper.DataSourceAttributes, *PropertyError) {
	objectAttributes := attrmapper.DataSourceAttributes{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, SchemaOpts{}, s.GlobalSchemaOpts)
		if err != nil {
			return nil, s.NewPropertyError(err, name)
		}

		attribute, propErr := pSchema.BuildDataSourceAttribute(name, s.GetComputability(name))
		if err != nil {
			return nil, propErr
		}

		objectAttributes = append(objectAttributes, attribute)
	}

	return objectAttributes, nil
}

func (s *OASSchema) BuildDataSourceAttribute(name string, computability schema.ComputedOptionalRequired) (attrmapper.DataSourceAttribute, *PropertyError) {
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
		return nil, s.NewPropertyError(fmt.Errorf("invalid schema type '%s'", s.Type), name)
	}
}

func (s *OASSchema) BuildProviderAttributes() (attrmapper.ProviderAttributes, *PropertyError) {
	objectAttributes := attrmapper.ProviderAttributes{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy, SchemaOpts{}, s.GlobalSchemaOpts)
		if err != nil {
			return nil, s.NewPropertyError(err, name)
		}

		attribute, propErr := pSchema.BuildProviderAttribute(name, s.GetOptionalOrRequired(name))
		if err != nil {
			return nil, propErr
		}

		objectAttributes = append(objectAttributes, attribute)
	}

	return objectAttributes, nil
}

func (s *OASSchema) BuildProviderAttribute(name string, optionalOrRequired schema.OptionalRequired) (attrmapper.ProviderAttribute, *PropertyError) {
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
		return nil, s.NewPropertyError(fmt.Errorf("invalid schema type '%s'", s.Type), name)
	}
}
