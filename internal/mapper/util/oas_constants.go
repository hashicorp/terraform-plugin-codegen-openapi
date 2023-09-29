// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package util

// Reference links:
//   - [JSON Schema - types]
//   - [JSON Schema - format]
//   - [JSON schema - custom format]
//   - [OAS - format]
//
// [JSON Schema - types]: https://json-schema.org/draft/2020-12/json-schema-core.html#name-instance-data-model
// [JSON Schema - format]: https://json-schema.org/draft/2020-12/json-schema-validation.html#name-defined-formats
// [JSON schema - custom format]: https://json-schema.org/draft/2020-12/json-schema-validation.html#name-custom-format-attributes
// [OAS - format]: https://spec.openapis.org/oas/latest.html#data-types
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

	OAS_param_path  = "path"
	OAS_param_query = "query"

	// Custom format for SetNested and Set attributes
	TF_format_set = "set"

	OAS_mediatype_json = "application/json"

	OAS_response_code_ok      = "200"
	OAS_response_code_created = "201"
)
