package schema

import (
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/util"
)

func (s *OASSchema) BuildResourceAttributes() (*[]ir.ResourceAttribute, error) {
	objectAttributes := []ir.ResourceAttribute{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy)
		if err != nil {
			return nil, err
		}

		attribute, err := pSchema.BuildResourceAttribute(name, s.GetBehavior(name))
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema - %w", name, err)
		}

		objectAttributes = append(objectAttributes, *attribute)
	}

	return &objectAttributes, nil
}

func (s *OASSchema) BuildResourceAttribute(name string, behavior ir.ComputedOptionalRequired) (*ir.ResourceAttribute, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringResource(name, behavior)
	case util.OAS_type_integer:
		return s.BuildIntegerResource(name, behavior)
	case util.OAS_type_number:
		return s.BuildNumberResource(name, behavior)
	case util.OAS_type_boolean:
		return s.BuildBoolResource(name, behavior)
	case util.OAS_type_array:
		return s.BuildListResource(name, behavior)
	case util.OAS_type_object:
		return s.BuildSingleNestedResource(name, behavior)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}

func (s *OASSchema) BuildDataSourceAttributes() (*[]ir.DataSourceAttribute, error) {
	objectAttributes := []ir.DataSourceAttribute{}

	// TODO: throw error if it's not an object?

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		pProxy := s.Schema.Properties[name]
		pSchema, err := BuildSchema(pProxy)
		if err != nil {
			return nil, err
		}

		attribute, err := pSchema.BuildDataSourceAttribute(name, s.GetBehavior(name))
		if err != nil {
			return nil, fmt.Errorf("failed to create object property '%s' schema - %w", name, err)
		}

		objectAttributes = append(objectAttributes, *attribute)
	}

	return &objectAttributes, nil
}

func (s *OASSchema) BuildDataSourceAttribute(name string, behavior ir.ComputedOptionalRequired) (*ir.DataSourceAttribute, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringDataSource(name, behavior)
	case util.OAS_type_integer:
		return s.BuildIntegerDataSource(name, behavior)
	case util.OAS_type_number:
		return s.BuildNumberDataSource(name, behavior)
	case util.OAS_type_boolean:
		return s.BuildBoolDataSource(name, behavior)
	case util.OAS_type_array:
		return s.BuildListDataSource(name, behavior)
	case util.OAS_type_object:
		return s.BuildSingleNestedDataSource(name, behavior)
	default:
		return nil, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}
