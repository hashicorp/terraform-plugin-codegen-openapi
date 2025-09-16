// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/util"

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

	s, err := buildSchemaProxy(proxy, globalOpts)
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
func buildSchemaProxy(proxy *base.SchemaProxy, globalOpts GlobalSchemaOpts) (*base.Schema, *SchemaError) {
	s, err := proxy.BuildSchema()
	if err != nil {
		return nil, SchemaErrorFromProxy(fmt.Errorf("failed to build schema proxy - %w", err), proxy)
	}

	// Check for discriminator patterns first - can work with oneOf, anyOf, allOf, or direct schemas
	if s.Discriminator != nil {
		// Prevent infinite recursion with depth limit
		if globalOpts.DiscriminatorDepth > 5 {
			// Skip discriminator processing at deep levels to prevent infinite recursion
			return s, nil
		}

		// Increment depth for recursive calls
		nextGlobalOpts := globalOpts
		nextGlobalOpts.DiscriminatorDepth++

		flattenedSchema, err := flattenDiscriminatorSchema(s, proxy, nextGlobalOpts)
		if err != nil {
			return nil, err
		}
		return flattenedSchema, nil
	}

	// If there are no schema composition keywords, return the schema
	if len(s.AllOf) == 0 && len(s.AnyOf) == 0 && len(s.OneOf) == 0 {
		return s, nil
	}

	if len(s.AnyOf) > 0 {
		if len(s.AnyOf) == 2 {
			schema, err := getMultiTypeSchema(s.AnyOf[0], s.AnyOf[1], globalOpts)
			if err != nil {
				return nil, err
			}

			return schema, nil
		}

		// Dynamic type currently not supported
		return nil, SchemaErrorFromNode(fmt.Errorf("found %d anyOf subschema(s), schema composition is currently not supported", len(s.AnyOf)), s, AnyOf)
	}

	if len(s.OneOf) > 0 {
		if len(s.OneOf) == 2 {
			schema, err := getMultiTypeSchema(s.OneOf[0], s.OneOf[1], globalOpts)
			if err != nil {
				return nil, err
			}

			return schema, nil
		}

		// Dynamic type currently not supported
		return nil, SchemaErrorFromNode(fmt.Errorf("found %d oneOf subschema(s), schema composition is currently not supported", len(s.OneOf)), s, OneOf)
	}

	// If there is just one allOf, we can use it as the schema
	if len(s.AllOf) == 1 {
		allOfSchema, err := buildSchemaProxy(s.AllOf[0], globalOpts)
		if err != nil {
			return nil, err
		}

		// Override the description w/ the parent if populated
		if s.Description != "" {
			allOfSchema.Description = s.Description
		}

		return allOfSchema, nil
	}

	// Handle multiple allOf schemas by merging their properties
	if len(s.AllOf) > 1 {
		composedSchema, err := composeAllOfSchemas(s, globalOpts)
		if err != nil {
			return nil, err
		}
		return composedSchema, nil
	}

	// No schema composition keywords found
	return nil, SchemaErrorFromNode(fmt.Errorf("found %d allOf subschema(s), schema composition is currently not supported", len(s.AllOf)), s, AllOf)
}

// getMultiTypeSchema will check the types of both schemas provided and will return the non-null schema. If a null schema type is not
// detected, an error will be returned as multi-types are not supported
func getMultiTypeSchema(proxyOne *base.SchemaProxy, proxyTwo *base.SchemaProxy, globalOpts GlobalSchemaOpts) (*base.Schema, *SchemaError) {
	firstSchema, err := buildSchemaProxy(proxyOne, globalOpts)
	if err != nil {
		return nil, err
	}

	secondSchema, err := buildSchemaProxy(proxyTwo, globalOpts)
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
		// https://github.com/starburstdata/terraform-plugin-codegen-openapi/issues/79
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

// composeAllOfSchemas merges multiple allOf schemas into a single composed schema
// This is essential for discriminator patterns which use allOf to compose base + variant schemas
func composeAllOfSchemas(baseSchema *base.Schema, globalOpts GlobalSchemaOpts) (*base.Schema, *SchemaError) {
	// Start with the base schema structure
	composedSchema := &base.Schema{
		Type:        baseSchema.Type,
		Format:      baseSchema.Format,
		Description: baseSchema.Description,
		Properties:  orderedmap.New[string, *base.SchemaProxy](),
		Required:    []string{},
		ParentProxy: baseSchema.ParentProxy,
	}

	// Copy base schema properties first
	if baseSchema.Properties != nil {
		for pair := range orderedmap.Iterate(context.TODO(), baseSchema.Properties) {
			composedSchema.Properties.Set(pair.Key(), pair.Value())
		}
	}
	composedSchema.Required = append(composedSchema.Required, baseSchema.Required...)

	// Process each allOf sub-schema and merge their properties
	for _, allOfProxy := range baseSchema.AllOf {
		allOfSchema, err := buildSchemaProxy(allOfProxy, globalOpts)
		if err != nil {
			return nil, err
		}

		// Merge properties from this allOf sub-schema
		if allOfSchema.Properties != nil {
			for pair := range orderedmap.Iterate(context.TODO(), allOfSchema.Properties) {
				key := pair.Key()
				value := pair.Value()
				// Later schemas override earlier ones (right-to-left merge)
				composedSchema.Properties.Set(key, value)
			}
		}

		// Merge required fields
		composedSchema.Required = append(composedSchema.Required, allOfSchema.Required...)

		// Use the type and format from the first concrete schema that has them
		if composedSchema.Type == nil && allOfSchema.Type != nil {
			composedSchema.Type = allOfSchema.Type
		}
		if composedSchema.Format == "" && allOfSchema.Format != "" {
			composedSchema.Format = allOfSchema.Format
		}

		// Use description from the last schema that has one
		if allOfSchema.Description != "" {
			composedSchema.Description = allOfSchema.Description
		}
	}

	return composedSchema, nil
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

// flattenDiscriminatorSchema handles discriminator patterns by flattening all properties from the base schema
// and all discriminator sub-schemas into a single comprehensive schema. This ensures all possible fields
// are included in the generated Terraform resource schema.
//
// Discriminators can work with oneOf, anyOf, allOf, or direct mapping patterns.
func flattenDiscriminatorSchema(baseSchema *base.Schema, baseProxy *base.SchemaProxy, globalOpts GlobalSchemaOpts) (*base.Schema, *SchemaError) {
	// Start with a copy of the base schema
	flattenedSchema := &base.Schema{
		Type:        baseSchema.Type,
		Format:      baseSchema.Format,
		Description: baseSchema.Description,
		Properties:  orderedmap.New[string, *base.SchemaProxy](),
		Required:    []string{},
		ParentProxy: baseProxy,
	}

	// Copy all base schema properties
	if baseSchema.Properties != nil {
		for pair := range orderedmap.Iterate(context.TODO(), baseSchema.Properties) {
			flattenedSchema.Properties.Set(pair.Key(), pair.Value())
		}
	}

	// Copy required fields from base schema
	flattenedSchema.Required = append(flattenedSchema.Required, baseSchema.Required...)

	// Collect all sub-schemas from different composition patterns
	var subSchemaProxies []*base.SchemaProxy

	// Handle oneOf discriminator pattern
	if len(baseSchema.OneOf) > 0 {
		subSchemaProxies = append(subSchemaProxies, baseSchema.OneOf...)
	}

	// Handle anyOf discriminator pattern
	if len(baseSchema.AnyOf) > 0 {
		subSchemaProxies = append(subSchemaProxies, baseSchema.AnyOf...)
	}

	// Handle allOf discriminator pattern (inheritance)
	if len(baseSchema.AllOf) > 0 {
		subSchemaProxies = append(subSchemaProxies, baseSchema.AllOf...)
	}

	// Handle discriminator mapping if no composition keywords but mapping exists
	if len(subSchemaProxies) == 0 && baseSchema.Discriminator != nil && baseSchema.Discriminator.Mapping != nil {
		// Resolve external schema references from discriminator mapping
		mappedSchemas, err := resolveDiscriminatorMapping(baseSchema.Discriminator.Mapping, globalOpts.Document)
		if err != nil {
			// Log error but continue with base schema - don't fail the entire process
			return flattenedSchema, nil
		}
		subSchemaProxies = append(subSchemaProxies, mappedSchemas...)
	}

	// Process each discriminator sub-schema and merge their properties
	for _, subProxy := range subSchemaProxies {
		subSchema, err := buildSchemaProxy(subProxy, globalOpts)
		if err != nil {
			// Log warning but continue - don't fail the entire process
			continue
		}

		// Merge properties from sub-schema
		if subSchema.Properties != nil {
			for pair := range orderedmap.Iterate(context.TODO(), subSchema.Properties) {
				key := pair.Key()
				value := pair.Value()
				// Only add if property doesn't already exist (base schema takes precedence)
				if _, exists := flattenedSchema.Properties.Get(key); !exists {
					flattenedSchema.Properties.Set(key, value)
				}
			}
		}

		// Merge required fields from sub-schema
		// Note: In discriminator patterns, sub-schema required fields become optional in the flattened schema
		// since they're only required for specific discriminator values
		for _, requiredField := range subSchema.Required {
			// Check if already in required list
			found := false
			for _, existing := range flattenedSchema.Required {
				if existing == requiredField {
					found = true
					break
				}
			}
			if !found {
				// For discriminator sub-schemas, we don't make these fields required in the flattened schema
				// They will be conditionally required based on discriminator value at runtime
			}
		}
	}

	return flattenedSchema, nil
}

// resolveDiscriminatorMapping resolves external schema references from a discriminator mapping
// and returns a slice of schema proxies for all mapped schemas
func resolveDiscriminatorMapping(mapping *orderedmap.Map[string, string], document *high.Document) ([]*base.SchemaProxy, *SchemaError) {
	var resolvedProxies []*base.SchemaProxy

	// Iterate through all discriminator mapping entries
	for pair := range orderedmap.Iterate(context.TODO(), mapping) {
		// discriminatorValue := pair.Key()   // e.g., "galaxy", "glue", "hive"
		schemaRef := pair.Value() // e.g., "#/components/schemas/S3CatalogGalaxyMetastore"

		// Resolve the schema reference using the document context
		resolvedProxy, err := resolveSchemaReference(schemaRef, document)
		if err != nil {
			// Continue with other mappings if one fails
			continue
		}

		resolvedProxies = append(resolvedProxies, resolvedProxy)
	}

	return resolvedProxies, nil
}

// resolveSchemaReference resolves a schema reference like "#/components/schemas/SchemaName"
// and returns the corresponding schema proxy
func resolveSchemaReference(reference string, document *high.Document) (*base.SchemaProxy, error) {
	// Check if this is a components/schemas reference
	if !strings.HasPrefix(reference, "#/components/schemas/") {
		return nil, fmt.Errorf("unsupported reference format: %s", reference)
	}

	// Extract the schema name from the reference
	schemaName := strings.TrimPrefix(reference, "#/components/schemas/")

	// Check if document is provided
	if document == nil {
		return nil, fmt.Errorf("cannot resolve reference %s: no document provided", reference)
	}

	// Look up the schema in the document's components
	if document.Components == nil || document.Components.Schemas == nil {
		return nil, fmt.Errorf("cannot resolve reference %s: no components/schemas in document", reference)
	}

	schemaProxy, exists := document.Components.Schemas.Get(schemaName)
	if !exists {
		return nil, fmt.Errorf("schema not found: %s", schemaName)
	}

	return schemaProxy, nil
}
