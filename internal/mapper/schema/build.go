package schema

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var ErrMultiTypeSchema = errors.New("unsupported multi-type schema, attribute cannot be created")
var ErrSchemaNotFound = errors.New("no compatible schema found")

// BuildSchemaFromRequest will extract and build the schema from the request body of an operation
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func BuildSchemaFromRequest(op *high.Operation) (*OASSchema, error) {
	if op == nil || op.RequestBody == nil || len(op.RequestBody.Content) == 0 {
		return nil, ErrSchemaNotFound
	}

	return getSchemaFromMediaType(op.RequestBody.Content)
}

// BuildSchemaFromResponse will extract and build the schema from the response body of an operation
//   - Response codes of 200 and then 201 will be prioritized, then will continue to the next available 2xx code
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func BuildSchemaFromResponse(op *high.Operation) (*OASSchema, error) {
	if op == nil || op.Responses == nil || len(op.Responses.Codes) == 0 {
		return nil, ErrSchemaNotFound
	}

	okResponse, ok := op.Responses.Codes[util.OAS_response_code_ok]
	if ok {
		return getSchemaFromMediaType(okResponse.Content)
	}

	createdResponse, ok := op.Responses.Codes[util.OAS_response_code_created]
	if ok {
		return getSchemaFromMediaType(createdResponse.Content)
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
			return getSchemaFromMediaType(responseCode.Content)
		}
	}

	return nil, ErrSchemaNotFound
}

func getSchemaFromMediaType(mediaTypes map[string]*high.MediaType) (*OASSchema, error) {
	// TODO: we might consider vendored JSON media types and maybe do a "contains" check instead of strict equality
	jsonMediaType, ok := mediaTypes[util.OAS_mediatype_json]
	if ok && jsonMediaType.Schema != nil {
		s, err := BuildSchema(jsonMediaType.Schema)
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
			s, err := BuildSchema(mediaType.Schema)
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
func BuildSchema(proxy *base.SchemaProxy) (*OASSchema, error) {
	resp := OASSchema{}

	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema proxy - %w", err)
	}

	resp.original = schema
	resp.Schema = schema

	if len(schema.AnyOf) == 2 {
		schema, err := getNullableSchema(schema.AnyOf[0], schema.AnyOf[1])
		if err != nil {
			return nil, err
		}

		resp.Schema = schema
	}

	if len(schema.OneOf) == 2 {
		schema, err := getNullableSchema(schema.OneOf[0], schema.OneOf[1])
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
		return "", errors.New("property does not have a type, attribute cannot be created")
	case 1:
		return schema.Type[0], nil
	case 2:
		if schema.Type[0] == util.OAS_type_null {
			return schema.Type[1], nil
		} else if schema.Type[1] == util.OAS_type_null {
			return schema.Type[0], nil
		}
	}

	return "", fmt.Errorf("%w - %v", ErrMultiTypeSchema, schema.Type)
}

// getNullableSchema will check the types of both schemas provided and will return the non-null schema. If a null schema type is not
// detected, an error will be returned as multi-types are not supported
func getNullableSchema(proxyOne *base.SchemaProxy, proxyTwo *base.SchemaProxy) (*base.Schema, error) {
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

	if firstType == util.OAS_type_null {
		return secondSchema, nil
	} else if secondType == util.OAS_type_null {
		return firstSchema, nil
	}

	return nil, fmt.Errorf("%w - %s %s", ErrMultiTypeSchema, firstType, secondType)
}
