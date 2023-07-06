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
