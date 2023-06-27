// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type OASSchema struct {
	Type   string
	Format string
	Schema *base.Schema

	GlobalSchemaOpts GlobalSchemaOpts
	SchemaOpts       SchemaOpts

	// TODO: Should we export the original schema? In the case of nullable schemas, we may want to consider using the original schema's description/annotations/etc.
	original *base.Schema
}

// GlobalSchemaOpts is passed recursively through built OASSchema structs. This is used for options that need to control
// the entire of a schema and it's potential nested schemas, like overriding computability. (Required, Optional, Computed)
type GlobalSchemaOpts struct {
	// OverrideComputability will set all attribute and nested attribute `ComputedOptionalRequired` fields
	// to this value. This ensures that an optional attribute from a higher precedence operation, such as a
	// create request for a resource, does not become required from a lower precedence operation, such as an
	// read response for a resource.
	OverrideComputability schema.ComputedOptionalRequired
}

// SchemaOpts is NOT passed recursively through built OASSchema structs, and will only be available to the top level schema. This is used
// for options that need to control just the top level schema, like overriding descriptions.
type SchemaOpts struct {
	// OverrideDescription will set the attribute description to this field if populated, otherwise the attribute description
	// will be set to the description field of the `schema`.
	OverrideDescription string
}

// IsMap will perform a type assertion on the `additionalProperties` field to determine if a map type
// is appropriate. See: https://json-schema.org/understanding-json-schema/reference/object.html#additional-properties
func (s *OASSchema) IsMap() bool {
	_, isMap := s.Schema.AdditionalProperties.(*base.SchemaProxy)
	return isMap
}

func (s *OASSchema) GetDescription() *string {
	if s.SchemaOpts.OverrideDescription != "" {
		return &s.SchemaOpts.OverrideDescription
	}

	// TODO: potentially use original description for nullable types?
	return &s.Schema.Description
}

func (s *OASSchema) IsSensitive() *bool {
	isSensitive := s.Format == util.OAS_format_password

	return &isSensitive
}

func (s *OASSchema) GetComputability(name string) schema.ComputedOptionalRequired {
	if s.GlobalSchemaOpts.OverrideComputability != "" {
		return s.GlobalSchemaOpts.OverrideComputability
	}

	for _, prop := range s.Schema.Required {
		if name == prop {
			return schema.Required
		}
	}

	return schema.ComputedOptional
}
