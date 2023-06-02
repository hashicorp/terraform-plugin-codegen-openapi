// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

// JSON schema base types: https://json-schema.org/draft/2020-12/json-schema-core.html#name-instance-data-model
// JSON schema available format types: https://json-schema.org/draft/2020-12/json-schema-validation.html#name-defined-formats
// OAS available format types: https://spec.openapis.org/oas/latest.html#data-types
//
// JSON schema Custom formats: https://json-schema.org/draft/2020-12/json-schema-validation.html#name-custom-format-attributes
const (
	OAS_type_string  = "string"
	OAS_type_integer = "integer"
	OAS_type_number  = "number"
	OAS_type_boolean = "boolean"
	OAS_type_array   = "array"
	OAS_type_object  = "object"
	OAS_type_null    = "null"

	OAS_format_double   = "double"
	OAS_format_float    = "float"
	OAS_format_password = "password"

	OAS_mediatype_json = "application/json"

	OAS_response_code_ok      = "200"
	OAS_response_code_created = "201"
)
