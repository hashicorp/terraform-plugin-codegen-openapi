package schema

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"

	"github.com/pb33f/libopenapi/datamodel/high/base"
)

type OASSchema struct {
	Type   string
	Format string
	Schema *base.Schema

	// TODO: Should we export the original schema? In the case of nullable schemas, we may want to consider using the original schema's description/annotations/etc.
	original *base.Schema
}

func (s *OASSchema) GetDescription() *string {
	// TODO: potentially use original description for nullable types?
	return &s.Schema.Description
}

func (s *OASSchema) IsSensitive() *bool {
	isSensitive := s.Format == util.OAS_format_password

	return &isSensitive
}

func (s *OASSchema) GetBehavior(name string) ir.ComputedOptionalRequired {
	for _, prop := range s.Schema.Required {
		if name == prop {
			return ir.Required
		}
	}

	return ir.ComputedOptional
}
