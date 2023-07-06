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
	// Float64ValidatorPackage is the name of the float64 validation package in
	// the framework validators module.
	Float64ValidatorPackage = "float64validator"
)

var (
	// Float64ValidatorCodeImport is a single allocation of the framework
	// validators module float64validator package import.
	Float64ValidatorCodeImport code.Import = CodeImport(Float64ValidatorPackage)
)

// Float64ValidatorOneOf returns a custom validator mapped to the Float64validator
// package OneOf function. If the values are nil or empty, nil is returned.
func Float64ValidatorOneOf(values []float64) *schema.CustomValidator {
	if len(values) == 0 {
		return nil
	}

	var schemaDefinition strings.Builder

	schemaDefinition.WriteString(Float64ValidatorPackage)
	schemaDefinition.WriteString(".OneOf(\n")

	for _, value := range values {
		schemaDefinition.WriteString(strconv.FormatFloat(value, 'f', -1, 64) + ",\n")
	}

	schemaDefinition.WriteString(")")

	return &schema.CustomValidator{
		Imports: []code.Import{
			Float64ValidatorCodeImport,
		},
		SchemaDefinition: schemaDefinition.String(),
	}
}
