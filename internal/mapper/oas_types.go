package mapper

// JSON schema base types: https://json-schema.org/draft/2020-12/json-schema-core.html#name-instance-data-model
// JSON schema available format types: https://json-schema.org/draft/2020-12/json-schema-validation.html#name-defined-formats
// OAS available format types: https://spec.openapis.org/oas/latest.html#data-types
//
// JSON schema Custom formats: https://json-schema.org/draft/2020-12/json-schema-validation.html#name-custom-format-attributes
const (
	oas_type_string  = "string"
	oas_type_integer = "integer"
	oas_type_number  = "number"
	oas_type_boolean = "boolean"
	oas_type_array   = "array"
	oas_type_object  = "object"

	oas_format_int64    = "int64"
	oas_format_double   = "double"
	oas_format_password = "password"
)
