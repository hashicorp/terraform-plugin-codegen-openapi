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
	// Int64ValidatorPackage is the name of the int64 validation package in
	// the framework validators module.
	Int64ValidatorPackage = "int64validator"
)

var (
	// Int64ValidatorCodeImport is a single allocation of the framework
	// validators module int64validator package import.
	Int64ValidatorCodeImport code.Import = CodeImport(Int64ValidatorPackage)
)

// Int64ValidatorAtLeast returns a custom validator mapped to the
// int64validator package AtLeast function.
func Int64ValidatorAtLeast(min int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(Int64ValidatorPackage)
	schemaDefinition.WriteString(".AtLeast(")
	schemaDefinition.WriteString(strconv.FormatInt(min, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			Int64ValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// Int64ValidatorAtMost returns a custom validator mapped to the
// int64validator package AtMost function.
func Int64ValidatorAtMost(max int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(Int64ValidatorPackage)
	schemaDefinition.WriteString(".AtMost(")
	schemaDefinition.WriteString(strconv.FormatInt(max, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			Int64ValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// Int64ValidatorBetween returns a custom validator mapped to the
// int64validator package Between function.
func Int64ValidatorBetween(min, max int64) *schema.CustomValidator {
	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(Int64ValidatorPackage)
	schemaDefinition.WriteString(".Between(")
	schemaDefinition.WriteString(strconv.FormatInt(min, 10))
	schemaDefinition.WriteString(", ")
	schemaDefinition.WriteString(strconv.FormatInt(max, 10))
	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			Int64ValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}

// Int64ValidatorOneOf returns a custom validator mapped to the int64validator
// package OneOf function. If the values are nil or empty, nil is returned.
func Int64ValidatorOneOf(values []int64) *schema.CustomValidator {
	if len(values) == 0 {
		return nil
	}

	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(Int64ValidatorPackage)
	schemaDefinition.WriteString(".OneOf(\n")

	for _, value := range values {
		schemaDefinition.WriteString(strconv.FormatInt(value, 10) + ",\n")
	}

	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			Int64ValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}
