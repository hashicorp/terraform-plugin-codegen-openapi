// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func (s *OASSchema) BuildElementType() (schema.ElementType, error) {
	switch s.Type {
	case util.OAS_type_string:
		return s.BuildStringElementType()
	case util.OAS_type_integer:
		return s.BuildIntegerElementType()
	case util.OAS_type_number:
		return s.BuildNumberElementType()
	case util.OAS_type_boolean:
		return s.BuildBoolElementType()
	case util.OAS_type_array:
		return s.BuildCollectionElementType()
	case util.OAS_type_object:
		return s.BuildObjectElementType()

	default:
		return schema.ElementType{}, fmt.Errorf("invalid schema type '%s'", s.Type)
	}
}
