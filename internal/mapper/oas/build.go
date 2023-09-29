// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var ErrMultiTypeSchema = errors.New("unsupported multi-type, attribute cannot be created")
var ErrSchemaNotFound = errors.New("no compatible schema found")

// BuildSchemaFromRequest will extract and build the schema from the request body of an operation
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func BuildSchemaFromRequest(op *high.Operation, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	if op == nil || op.RequestBody == nil || len(op.RequestBody.Content) == 0 {
		return nil, ErrSchemaNotFound
	}

	return getSchemaFromMediaType(op.RequestBody.Content, schemaOpts, globalOpts)
}

// BuildSchemaFromResponse will extract and build the schema from the response body of an operation
//   - Response codes of 200 and then 201 will be prioritized, then will continue to the next available 2xx code
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func BuildSchemaFromResponse(op *high.Operation, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	if op == nil || op.Responses == nil || len(op.Responses.Codes) == 0 {
		return nil, ErrSchemaNotFound
	}

	okResponse, ok := op.Responses.Codes[util.OAS_response_code_ok]
	if ok {
		return getSchemaFromMediaType(okResponse.Content, schemaOpts, globalOpts)
	}

	createdResponse, ok := op.Responses.Codes[util.OAS_response_code_created]
	if ok {
		return getSchemaFromMediaType(createdResponse.Content, schemaOpts, globalOpts)
	}

	// Guarantee the order of processing
	codes := util.SortedKeys(op.Responses.Codes)
	for _, code := range codes {
		responseCode := op.Responses.Codes[code]
		statusCode, err := strconv.Atoi(code)
		if err != nil {
			continue
		}

		if statusCode >= 200 && statusCode <= 299 {
			return getSchemaFromMediaType(responseCode.Content, schemaOpts, globalOpts)
		}
	}

	return nil, ErrSchemaNotFound
}

func getSchemaFromMediaType(mediaTypes map[string]*high.MediaType, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	jsonMediaType, ok := mediaTypes[util.OAS_mediatype_json]
	if ok && jsonMediaType.Schema != nil {
		s, err := BuildSchema(jsonMediaType.Schema, schemaOpts, globalOpts)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	// Guarantee the order of processing
	mediaTypeKeys := util.SortedKeys(mediaTypes)
	for _, key := range mediaTypeKeys {
		mediaType := mediaTypes[key]
		if mediaType.Schema != nil {
			s, err := BuildSchema(mediaType.Schema, schemaOpts, globalOpts)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}

	return nil, ErrSchemaNotFound
}

// BuildSchema will build a schema from a schema proxy. It can also handle nullable schemas/types,
// implemented with oneOf/anyOf OAS keywords or an array on the "type" property
func BuildSchema(proxy *base.SchemaProxy, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	resp := OASSchema{}

	s, err := proxy.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema proxy - %w", err)
	}

	resp.SchemaOpts = schemaOpts
	resp.GlobalSchemaOpts = globalOpts

	resp.original = s
	resp.Schema = s

	if len(resp.Schema.AllOf) > 0 {
		// If there is just one allOf, we can use it as the schema
		if len(resp.Schema.AllOf) == 1 {
			schema, err := resp.Schema.AllOf[0].BuildSchema()
			if err != nil {
				return nil, fmt.Errorf("failed to build allOf[0] schema proxy - %w", err)
			}

			// Override the description w/ the parent if populated
			if resp.Schema.Description != "" {
				schema.Description = resp.Schema.Description
			}

			resp.Schema = schema
		}
	}

	if len(resp.Schema.AnyOf) == 2 {
		schema, err := getMultiTypeSchema(resp.Schema.AnyOf[0], resp.Schema.AnyOf[1])
		if err != nil {
			return nil, err
		}

		resp.Schema = schema
	}

	if len(resp.Schema.OneOf) == 2 {
		schema, err := getMultiTypeSchema(resp.Schema.OneOf[0], resp.Schema.OneOf[1])
		if err != nil {
			return nil, err
		}

		resp.Schema = schema
	}

	oasType, err := retrieveType(resp.Schema)
	if err != nil {
		return nil, err
	}

	resp.Type = oasType
	resp.Format = resp.Schema.Format

	return &resp, nil
}

// retrieveType will return the JSON schema type. Support for multi-types is restricted to combinations of "null" and another type, i.e. ["null", "string"]
func retrieveType(schema *base.Schema) (string, error) {
	switch len(schema.Type) {
	case 0:
		return "", errors.New("property does not have a 'type' or supported allOf, oneOf, anyOf constraint - attribute cannot be created")
	case 1:
		return schema.Type[0], nil
	case 2:
		// Check for null type, if found, return the other type
		if schema.Type[0] == util.OAS_type_null {
			return schema.Type[1], nil
		} else if schema.Type[1] == util.OAS_type_null {
			return schema.Type[0], nil
		}

		// Check for string type, if the other type can be represented as a string, return the string type
		if schema.Type[0] == util.OAS_type_string && isStringableType(schema.Type[1]) {
			return schema.Type[0], nil
		} else if schema.Type[1] == util.OAS_type_string && isStringableType(schema.Type[0]) {
			return schema.Type[1], nil
		}
	}

	return "", fmt.Errorf("%v - %w", schema.Type, ErrMultiTypeSchema)
}

// getMultiTypeSchema will check the types of both schemas provided and will return the non-null schema. If a null schema type is not
// detected, an error will be returned as multi-types are not supported
func getMultiTypeSchema(proxyOne *base.SchemaProxy, proxyTwo *base.SchemaProxy) (*base.Schema, error) {
	firstSchema, err := proxyOne.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema proxy - %w", err)
	}

	secondSchema, err := proxyTwo.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema proxy - %w", err)
	}

	firstType, err := retrieveType(firstSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve schema type - %w", err)
	}

	secondType, err := retrieveType(secondSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve schema type - %w", err)
	}

	// Check for null type, if found, return the other type
	if firstType == util.OAS_type_null {
		return secondSchema, nil
	} else if secondType == util.OAS_type_null {
		return firstSchema, nil
	}

	// Check for string type, if the other type can be represented as a string, return the string type
	if firstType == util.OAS_type_string && isStringableType(secondType) {
		return firstSchema, nil
	} else if secondType == util.OAS_type_string && isStringableType(firstType) {
		return secondSchema, nil
	}

	return nil, fmt.Errorf("[%s, %s] - %w", firstType, secondType, ErrMultiTypeSchema)
}

func isStringableType(t string) bool {
	switch t {
	case util.OAS_type_string:
		return true
	case util.OAS_type_integer:
		return true
	case util.OAS_type_number:
		return true
	case util.OAS_type_boolean:
		return true
	default:
		return false
	}
}
