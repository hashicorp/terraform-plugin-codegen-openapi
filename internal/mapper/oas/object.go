// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildObjectElementType() (schema.ElementType, *SchemaError) {
	objectElemTypes := []schema.ObjectAttributeType{}

	// Guarantee the order of processing
	propertyNames := util.SortedKeys(s.Schema.Properties)
	for _, name := range propertyNames {

		if s.IsPropertyIgnored(name) {
			continue
		}

		pProxy := s.Schema.Properties[name]
		schemaOpts := SchemaOpts{
			Ignores: s.GetIgnoresForNested(name),
		}

		pSchema, err := BuildSchema(pProxy, schemaOpts, s.GlobalSchemaOpts)
		if err != nil {
			return schema.ElementType{}, s.NestSchemaError(err, name)
		}

		elemType, err := pSchema.BuildElementType()
		if err != nil {
			return schema.ElementType{}, s.NestSchemaError(err, name)
		}

		objectElemTypes = append(objectElemTypes, util.CreateObjectAttributeType(name, elemType))
	}

	return schema.ElementType{
		Object: &schema.ObjectType{
			AttributeTypes: objectElemTypes,
		},
	}, nil
}
