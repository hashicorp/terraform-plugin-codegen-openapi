// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"context"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/orderedmap"
)

func (s *OASSchema) BuildObjectElementType() (schema.ElementType, *SchemaError) {
	objectElemTypes := []schema.ObjectAttributeType{}

	sortedProperties := orderedmap.SortAlpha(s.Schema.Properties)
	for pair := range orderedmap.Iterate(context.TODO(), sortedProperties) {
		name := pair.Key()

		if s.IsPropertyIgnored(name) {
			continue
		}

		pProxy := pair.Value()
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
