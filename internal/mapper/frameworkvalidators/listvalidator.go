// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkvalidators

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

const (
	// ListValidatorPackage is the name of the list validation package in
	// the framework validators module.
	ListValidatorPackage = "listvalidator"
)

var (
	// ListValidatorCodeImport is a single allocation of the framework
	// validators module listvalidator package import.
	ListValidatorCodeImport code.Import = CodeImport(ListValidatorPackage)
)

// ListValidatorSizeAtLeast returns a custom validator mapped to the
// listvalidator package SizeAtLeast function.
func ListValidatorSizeAtLeast(minimum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(ListValidatorPackage)
	schemaDefinition.WriteString(".SizeAtLeast(")
	schemaDefinition.WriteString(strconv.FormatInt(minimum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			ListValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// ListValidatorSizeAtMost returns a custom validator mapped to the
// listvalidator package SizeAtMost function.
func ListValidatorSizeAtMost(maximum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(ListValidatorPackage)
	schemaDefinition.WriteString(".SizeAtMost(")
	schemaDefinition.WriteString(strconv.FormatInt(maximum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			ListValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// ListValidatorSizeBetween returns a custom validator mapped to the
// listvalidator package SizeBetween function.
func ListValidatorSizeBetween(minimum, maximum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(ListValidatorPackage)
	schemaDefinition.WriteString(".SizeBetween(")
	schemaDefinition.WriteString(strconv.FormatInt(minimum, 10))
	schemaDefinition.WriteString(", ")
	schemaDefinition.WriteString(strconv.FormatInt(maximum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			ListValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// ListValidatorUniqueValues returns a custom validator mapped to the
// listvalidator package UniqueValues function.
func ListValidatorUniqueValues() *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(ListValidatorPackage)
	schemaDefinition.WriteString(".UniqueValues()")

	return &schema.CustomValidator{
		Imports: []code.Import{
			ListValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}
