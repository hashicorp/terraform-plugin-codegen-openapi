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
	// SetValidatorPackage is the name of the set validation package in
	// the framework validators module.
	SetValidatorPackage = "setvalidator"
)

var (
	// SetValidatorCodeImport is a single allocation of the framework
	// validators module setvalidator package import.
	SetValidatorCodeImport code.Import = CodeImport(SetValidatorPackage)
)

// SetValidatorSizeAtLeast returns a custom validator mapped to the
// Setvalidator package SizeAtLeast function.
func SetValidatorSizeAtLeast(minimum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(SetValidatorPackage)
	schemaDefinition.WriteString(".SizeAtLeast(")
	schemaDefinition.WriteString(strconv.FormatInt(minimum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			SetValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// SetValidatorSizeAtMost returns a custom validator mapped to the
// Setvalidator package SizeAtMost function.
func SetValidatorSizeAtMost(maximum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(SetValidatorPackage)
	schemaDefinition.WriteString(".SizeAtMost(")
	schemaDefinition.WriteString(strconv.FormatInt(maximum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			SetValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// SetValidatorSizeBetween returns a custom validator mapped to the
// Setvalidator package SizeBetween function.
func SetValidatorSizeBetween(minimum, maximum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(SetValidatorPackage)
	schemaDefinition.WriteString(".SizeBetween(")
	schemaDefinition.WriteString(strconv.FormatInt(minimum, 10))
	schemaDefinition.WriteString(", ")
	schemaDefinition.WriteString(strconv.FormatInt(maximum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			SetValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}
