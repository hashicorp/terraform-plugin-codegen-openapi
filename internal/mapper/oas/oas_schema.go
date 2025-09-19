// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"context"
	"strings"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

type OASSchema struct {
	Type   string
	Format string
	Schema *base.Schema

	GlobalSchemaOpts GlobalSchemaOpts
	SchemaOpts       SchemaOpts
}

// GlobalSchemaOpts is passed recursively through built OASSchema structs. This is used for options that need to control
// the entire of a schema and it's potential nested schemas, like overriding computability. (Required, Optional, Computed)
type GlobalSchemaOpts struct {
	// OverrideComputability will set all attribute and nested attribute `ComputedOptionalRequired` fields
	// to this value. This ensures that an optional attribute from a higher precedence operation, such as a
	// create request for a resource, does not become required from a lower precedence operation, such as an
	// read response for a resource.
	OverrideComputability schema.ComputedOptionalRequired

	// Document provides access to the full OpenAPI document for resolving external schema references
	// in discriminator mappings (e.g., "#/components/schemas/S3CatalogGalaxyMetastore")
	Document *high.Document

	// DiscriminatorDepth tracks recursion depth to prevent infinite loops in discriminator resolution
	DiscriminatorDepth int
}

// SchemaOpts is NOT passed recursively through built OASSchema structs, and will only be available to the top level schema. This is used
// for options that need to control just the top level schema, like overriding descriptions.
type SchemaOpts struct {
	// Ignores contains all potentially relevant ignores for a schema and it's potential nested schemas
	Ignores []string

	// OverrideDeprecationMessage will set the attribute deprecation message to
	// this field if populated, otherwise the attribute deprecation message will
	// be set to a default "This attribute is deprecated." message when the
	// deprecated property is enabled.
	OverrideDeprecationMessage string

	// OverrideDescription will set the attribute description to this field if populated, otherwise the attribute description
	// will be set to the description field of the `schema`.
	OverrideDescription string
}

// IsMap checks the `additionalProperties` field to determine if a map type is appropriate (refer to [JSON Schema - additionalProperties]).
//
// [JSON Schema - additionalProperties]: https://json-schema.org/understanding-json-schema/reference/object.html#additional-properties
func (s *OASSchema) IsMap() bool {
	return s.Schema.AdditionalProperties != nil && s.Schema.AdditionalProperties.IsA()
}

// SchemaErrorFromProperty is a helper function for creating an SchemaError struct for a property.
func (s *OASSchema) SchemaErrorFromProperty(err error, propName string) *SchemaError {
	return NewSchemaError(err, s.getPropertyLineNumber(propName), propName)
}

// NestSchemaError is a helper function for creating a nested SchemaError struct for a property.
func (s *OASSchema) NestSchemaError(err *SchemaError, propName string) *SchemaError {
	return err.NestedSchemaError(propName, s.getPropertyLineNumber(propName))
}

// getPropertyLineNumber looks in the low-level schema instance for line information. Returns 0 if not found.
func (s *OASSchema) getPropertyLineNumber(propName string) int {
	low := s.Schema.GoLow()
	if low == nil {
		return 0
	}

	// Check property nodes first for a line number
	for pair := range orderedmap.Iterate(context.TODO(), low.Properties.Value) {
		if pair.Key().Value == propName {
			return pair.Value().NodeLineNumber()
		}
	}

	// If it's not found in properties, default to the line number from the parent node
	if low.ParentProxy != nil && low.ParentProxy.GetValueNode() != nil {
		return low.ParentProxy.GetValueNode().Line
	}

	return 0
}

// GetDeprecationMessage returns a deprecation message if the deprecated
// property is enabled. It defaults the message to "This attribute is
// deprecated" unless the SchemaOpts.OverrideDeprecationMessage is set.
func (s *OASSchema) GetDeprecationMessage() *string {
	if s.Schema.Deprecated == nil || !(*s.Schema.Deprecated) {
		return nil
	}

	if s.SchemaOpts.OverrideDeprecationMessage != "" {
		return &s.SchemaOpts.OverrideDeprecationMessage
	}

	deprecationMessage := "This attribute is deprecated."

	return &deprecationMessage
}

func (s *OASSchema) GetDescription() *string {
	if s.SchemaOpts.OverrideDescription != "" {
		return &s.SchemaOpts.OverrideDescription
	}

	if s.Schema.Description == "" {
		return nil
	}

	return &s.Schema.Description
}

func (s *OASSchema) IsSensitive() *bool {
	isSensitive := s.Format == util.OAS_format_password

	if !isSensitive {
		return nil
	}

	return &isSensitive
}

// TODO: Figure out a better way to handle computability, since it differs with provider vs. datasource/resource
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

func (s *OASSchema) GetOptionalOrRequired(name string) schema.OptionalRequired {
	for _, prop := range s.Schema.Required {
		if name == prop {
			return schema.Required
		}
	}

	return schema.Optional
}

// IsPropertyIgnored checks if a property should be ignored
func (s *OASSchema) IsPropertyIgnored(name string) bool {
	for _, ignore := range s.SchemaOpts.Ignores {
		if name == ignore {
			return true
		}
	}
	return false
}

// GetIgnoresForNested is a helper function that will return all nested ignores for a property. If no ignores
// or nested ignores are found, returns an empty string slice.
func (s *OASSchema) GetIgnoresForNested(name string) []string {
	newIgnores := make([]string, 0)

	for _, ignore := range s.SchemaOpts.Ignores {
		ignoreParts := strings.Split(ignore, ".")

		if len(ignoreParts) > 1 && name == ignoreParts[0] {
			newIgnore := strings.Join(ignoreParts[1:], ".")

			if newIgnore != "" {
				newIgnores = append(newIgnores, newIgnore)
			}
		}
	}

	return newIgnores
}
