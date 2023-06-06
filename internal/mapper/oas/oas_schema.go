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

	// TODO: Should we export the original schema? In the case of nullable schemas, we may want to consider using the original schema's description/annotations/etc.
	original *base.Schema
}

// GlobalSchemaOpts is passed recursively through built OASSchema structs
type GlobalSchemaOpts struct {
	// OverrideComputability will default all attribute's and nested attribute's `ComputedOptionalRequired` field to this value
	OverrideComputability schema.ComputedOptionalRequired
}

func (s *OASSchema) GetDescription() *string {
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
