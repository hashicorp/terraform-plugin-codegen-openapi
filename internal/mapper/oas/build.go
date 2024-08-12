// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/util"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

var ErrMultiTypeSchema = errors.New("unsupported multi-type, attribute cannot be created")
var ErrSchemaNotFound = errors.New("no compatible schema found")

// BuildSchemaFromRequest will extract and build the schema from the request body of an operation
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func BuildSchemaFromRequest(op *high.Operation, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	if op == nil || op.RequestBody == nil || op.RequestBody.Content == nil || op.RequestBody.Content.Len() == 0 {
		return nil, ErrSchemaNotFound
	}

	return getSchemaFromMediaType(op.RequestBody.Content, schemaOpts, globalOpts)
}

// BuildSchemaFromResponse will extract and build the schema from the response body of an operation
//   - Response codes of 200 and then 201 will be prioritized, then will continue to the next available 2xx code
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func BuildSchemaFromResponse(op *high.Operation, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	if op == nil || op.Responses == nil || op.Responses.Codes == nil || op.Responses.Codes.Len() == 0 {
		return nil, ErrSchemaNotFound
	}

	okResponse, ok := op.Responses.Codes.Get(util.OAS_response_code_ok)
	if ok {
		return getSchemaFromMediaType(okResponse.Content, schemaOpts, globalOpts)
	}

	createdResponse, ok := op.Responses.Codes.Get(util.OAS_response_code_created)
	if ok {
		return getSchemaFromMediaType(createdResponse.Content, schemaOpts, globalOpts)
	}

	sortedCodes := orderedmap.SortAlpha(op.Responses.Codes)
	for pair := range orderedmap.Iterate(context.TODO(), sortedCodes) {
		responseCode := pair.Value()
		statusCode, err := strconv.Atoi(pair.Key())
		if err != nil {
			continue
		}

		if statusCode >= 200 && statusCode <= 299 {
			return getSchemaFromMediaType(responseCode.Content, schemaOpts, globalOpts)
		}
	}

	return nil, ErrSchemaNotFound
}

func getSchemaFromMediaType(mediaTypes *orderedmap.Map[string, *high.MediaType], schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, error) {
	if mediaTypes == nil {
		return nil, ErrSchemaNotFound
	}

	jsonMediaType, ok := mediaTypes.Get(util.OAS_mediatype_json)
	if ok && jsonMediaType.Schema != nil {
		s, err := BuildSchema(jsonMediaType.Schema, schemaOpts, globalOpts)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	sortedMediaTypes := orderedmap.SortAlpha(mediaTypes)
	for pair := range orderedmap.Iterate(context.TODO(), sortedMediaTypes) {
		mediaType := pair.Value()
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
func BuildSchema(proxy *base.SchemaProxy, schemaOpts SchemaOpts, globalOpts GlobalSchemaOpts) (*OASSchema, *SchemaError) {
	resp := OASSchema{}

	s, err := buildSchemaProxy(proxy)
	if err != nil {
		return nil, err
	}

	resp.SchemaOpts = schemaOpts
	resp.GlobalSchemaOpts = globalOpts
	resp.Schema = s

	oasType, err := retrieveType(resp.Schema)
	if err != nil {
		return nil, err
	}

	resp.Type = oasType
	resp.Format = resp.Schema.Format

	return &resp, nil
}

// buildSchemaProxy is a helper function that builds a schema proxy. If needed, it will recursively resolve a specific set of [schema composition] keywords:
//   - allOf: If len == 1, will resolve with that one item.
//   - anyOf: If len == 2, will resolve nullable or stringable types
//   - oneOf: If len == 2, will resolve nullable or stringable types
//
// # Any other combinations of allOf, anyOf, or oneOf will return a SchemaError
//
// [schema composition]: https://json-schema.org/understanding-json-schema/reference/combining
func buildSchemaProxy(proxy *base.SchemaProxy) (*base.Schema, *SchemaError) {
	s, err := proxy.BuildSchema()
	if err != nil {
		return nil, SchemaErrorFromProxy(fmt.Errorf("failed to build schema proxy - %w", err), proxy)
	}

	// If there are no schema composition keywords, return the schema
	if len(s.AllOf) == 0 && len(s.AnyOf) == 0 && len(s.OneOf) == 0 {
		return s, nil
	}

	if len(s.AnyOf) > 0 {
		if len(s.AnyOf) == 2 {
			schema, err := getMultiTypeSchema(s.AnyOf[0], s.AnyOf[1])
			if err == nil {
				return schema, nil
			}
		}

		schema, err := handleOneOfAnyOf(s.AnyOf)
		if err != nil {
			return nil, SchemaErrorFromNode(err, s, AnyOf)
		}

		return schema, nil
	}

	if len(s.OneOf) > 0 {
		if len(s.OneOf) == 2 {
			schema, err := getMultiTypeSchema(s.OneOf[0], s.OneOf[1])
			if err == nil {
				return schema, nil
			}
		}

		schema, err := handleOneOfAnyOf(s.OneOf)
		if err != nil {
			return nil, SchemaErrorFromNode(err, s, OneOf)
		}

		return schema, nil
	}

	// If there is just one allOf, we can use it as the schema
	if len(s.AllOf) == 1 {
		allOfSchema, err := buildSchemaProxy(s.AllOf[0])
		if err != nil {
			return nil, err
		}

		// Override the description w/ the parent if populated
		if s.Description != "" {
			allOfSchema.Description = s.Description
		}

		return allOfSchema, nil
	}

	compoundSchema, err := compoundAllOf(s)
	if err != nil {
		return nil, SchemaErrorFromNode(fmt.Errorf("%w, schema composition is currently not supported", err), s, AllOf)
	}

	return compoundSchema, nil
}

var oneOfAnyOfReg = regexp.MustCompile("[^a-zA-Z0-9_]+")

func handleOneOfAnyOf(s []*base.SchemaProxy) (*base.Schema, error) {
	compoundSchema := &base.Schema{Properties: orderedmap.New[string, *base.SchemaProxy]()}
	for i, schema := range s {
		schema, err := buildSchemaProxy(schema)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", i, err)
		}

		key := "field_" + strconv.Itoa(i)
		if schema.Title != "" {
			key = schema.Title
			key = strings.ToLower(key)
			key = strings.ReplaceAll(key, " ", "_")
			key = oneOfAnyOfReg.ReplaceAllString(key, "")
		}

		compoundSchema.Properties.Store(key, base.CreateSchemaProxy(schema))
	}

	return compoundSchema, nil
}

func compoundAllOf(s *base.Schema) (*base.Schema, error) {
	var commonType string
	var schemas []*base.Schema
	for i, schemaProxy := range s.AllOf {
		schema, err := buildSchemaProxy(schemaProxy)
		if err != nil {
			return nil, fmt.Errorf("%v: building subschema", i)
		}

		schemas = append(schemas, schema)

		typ, err := retrieveType(schema)
		if err != nil {
			return nil, fmt.Errorf("%v: %w", i, err)
		}

		if commonType == "" {
			commonType = typ
		} else if commonType != typ {
			return nil, fmt.Errorf("%v: different types, got %v, expected %v", i, typ, commonType)
		}
	}

	switch commonType {
	case util.OAS_type_object:
		return compoundAllOfObject(schemas)
	case util.OAS_type_string:
		return compoundAllOfString(schemas)
	default:
		return nil, fmt.Errorf("unhandled type: %v", commonType)
	}
}

func compoundAllOfObject(s []*base.Schema) (*base.Schema, error) {
	compoundSchema := &base.Schema{Properties: orderedmap.New[string, *base.SchemaProxy]()}
	for i, schema := range s {
		pair := schema.Properties.First()
		for pair != nil {
			pairSchema, err := buildSchemaProxy(pair.Value())
			if err != nil {
				return nil, fmt.Errorf("%v: %v: %w", i, pair.Key(), err)
			}

			compoundSchema.Properties.Store(pair.Key(), base.CreateSchemaProxy(pairSchema))

			pair = pair.Next()
		}
	}

	return compoundSchema, nil
}

func compoundAllOfString(s []*base.Schema) (*base.Schema, error) {
	// TODO: enum
	return &base.Schema{Type: []string{"string"}}, nil
}

// getMultiTypeSchema will check the types of both schemas provided and will return the non-null schema. If a null schema type is not
// detected, an error will be returned as multi-types are not supported
func getMultiTypeSchema(proxyOne *base.SchemaProxy, proxyTwo *base.SchemaProxy) (*base.Schema, *SchemaError) {
	firstSchema, err := buildSchemaProxy(proxyOne)
	if err != nil {
		return nil, err
	}

	secondSchema, err := buildSchemaProxy(proxyTwo)
	if err != nil {
		return nil, err
	}

	firstType, err := retrieveType(firstSchema)
	if err != nil {
		return nil, err
	}

	secondType, err := retrieveType(secondSchema)
	if err != nil {
		return nil, err
	}

	if firstType == secondType {
		s, err := compoundAllOfString([]*base.Schema{firstSchema, secondSchema})
		if err != nil {
			return nil, SchemaErrorFromNode(err, firstSchema, Type)
		}

		return s, nil
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

	return nil, SchemaErrorFromNode(fmt.Errorf("[%s %s] - %w", firstType, secondType, ErrMultiTypeSchema), firstSchema, Type)
}

// retrieveType will return the JSON schema type. Support for multi-types is restricted to combinations of "null" and another type, i.e. ["null", "string"]
func retrieveType(schema *base.Schema) (string, *SchemaError) {
	switch len(schema.Type) {
	case 0:
		// Properties are only valid applying to objects, it's possible tools might omit the type
		// https://github.com/raphaelfff/terraform-plugin-codegen-openapi/issues/79
		if schema.Properties != nil && schema.Properties.Len() > 0 {
			return util.OAS_type_object, nil
		}

		return "", SchemaErrorFromProxy(errors.New("no 'type' array or supported allOf, oneOf, anyOf constraint - attribute cannot be created"), schema.ParentProxy)
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

	return "", SchemaErrorFromNode(fmt.Errorf("%v - %w", schema.Type, ErrMultiTypeSchema), schema, Type)
}

func isStringableType(t string) bool {
	switch t {
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
