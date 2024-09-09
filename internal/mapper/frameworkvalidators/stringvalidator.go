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
	// StringValidatorPackage is the name of the string validation package in
	// the framework validators module.
	StringValidatorPackage = "stringvalidator"
)

var (
	// StringValidatorCodeImport is a single allocation of the framework
	// validators module stringvalidator package import.
	StringValidatorCodeImport code.Import = CodeImport(StringValidatorPackage)
)

// StringValidatorLengthAtLeast returns a custom validator mapped to the
// stringvalidator package LengthAtLeast function.
func StringValidatorLengthAtLeast(minimum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(StringValidatorPackage)
	schemaDefinition.WriteString(".LengthAtLeast(")
	schemaDefinition.WriteString(strconv.FormatInt(minimum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			StringValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// StringValidatorLengthAtMost returns a custom validator mapped to the
// stringvalidator package LengthAtMost function.
func StringValidatorLengthAtMost(maximum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(StringValidatorPackage)
	schemaDefinition.WriteString(".LengthAtMost(")
	schemaDefinition.WriteString(strconv.FormatInt(maximum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			StringValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// StringValidatorLengthBetween returns a custom validator mapped to the
// stringvalidator package LengthBetween function.
func StringValidatorLengthBetween(minimum, maximum int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(StringValidatorPackage)
	schemaDefinition.WriteString(".LengthBetween(")
	schemaDefinition.WriteString(strconv.FormatInt(minimum, 10))
	schemaDefinition.WriteString(", ")
	schemaDefinition.WriteString(strconv.FormatInt(maximum, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			StringValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// StringValidatorOneOf returns a custom validator mapped to the stringvalidator
// package OneOf function. If the values are nil or empty, nil is returned.
func StringValidatorOneOf(values []string) *schema.CustomValidator {
	if len(values) == 0 {
		return nil
	}

	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(StringValidatorPackage)
	schemaDefinition.WriteString(".OneOf(\n")

	for _, value := range values {
		schemaDefinition.WriteString(strconv.Quote(value) + ",\n")
	}

	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			StringValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// StringValidatorRegexMatches returns a custom validator mapped to the
// stringvalidator package RegexMatches function.
func StringValidatorRegexMatches(pattern, message string) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(StringValidatorPackage)
	schemaDefinition.WriteString(".RegexMatches(")
	schemaDefinition.WriteString("regexp.MustCompile(")
	schemaDefinition.WriteString(strconv.Quote(pattern))
	schemaDefinition.WriteString(")")
	schemaDefinition.WriteString(", ")
	schemaDefinition.WriteString(strconv.Quote(message))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			{
				Path: "regexp",
			},
			StringValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}
