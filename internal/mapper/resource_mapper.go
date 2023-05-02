package mapper

import (
	"errors"
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"log"
	"sort"
	"strconv"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var _ ResourceMapper = resourceMapper{}

var errMultiTypeSchema = errors.New("unsupported multi-type schema, attribute cannot be created")

type ResourceMapper interface {
	MapToIR(resources map[string]explorer.Resource) (*[]ir.Resource, error)
}

type resourceMapper struct{}

func NewResourceMapper() ResourceMapper {
	return resourceMapper{}
}

func (m resourceMapper) MapToIR(resources map[string]explorer.Resource) (*[]ir.Resource, error) {
	resourceSchemas := []ir.Resource{}

	// Guarantee the order of processing
	resourceNames := sortedKeys(resources)
	for _, name := range resourceNames {
		resource := resources[name]

		schema, err := generateResourceSchema(resource)
		if err != nil {
			log.Printf("[WARN] skipping '%s' resource schema: %s\n", name, err)
			continue
		}

		resourceSchemas = append(resourceSchemas, ir.Resource{
			Name:   name,
			Schema: *schema,
		})
	}

	return &resourceSchemas, nil
}

func generateResourceSchema(resource explorer.Resource) (*ir.ResourceSchema, error) {
	resourceSchema := &ir.ResourceSchema{
		Attributes: []ir.ResourceAttribute{},
	}

	// ***************
	// Create Request Body
	// ***************
	requestBodySchema := getRequestSchemaProxy(resource.CreateOp)
	if requestBodySchema == nil {
		return nil, errors.New("no request schema found!")
	}

	requestBodyAttributes, err := mapSchemaToObjectAttributes(requestBodySchema)
	if err != nil {
		return nil, err
	}

	// ***************
	// Create Response Body
	// ***************
	createResponseAttributes := &[]ir.ResourceAttribute{}
	createResponseSchema := getResponseSchemaProxy(resource.CreateOp)
	if createResponseSchema != nil {
		createResponseAttributes, err = mapSchemaToObjectAttributes(createResponseSchema)
		if err != nil {
			return nil, err
		}
	}

	// ***************
	// READ Response Body
	// ***************
	readResponseAttributes := &[]ir.ResourceAttribute{}
	responseBodySchema := getResponseSchemaProxy(resource.ReadOp)
	if responseBodySchema != nil {
		readResponseAttributes, err = mapSchemaToObjectAttributes(responseBodySchema)
		if err != nil {
			return nil, err
		}
	}

	// ***************
	// READ Parameters
	// ***************
	// TODO: Merging in root of schema, not sure a better way of handling
	parameterAttributes := []ir.ResourceAttribute{}
	for _, param := range getResponseParameters(resource.ReadOp) {
		// TODO: Filter specific "in" values? - https://spec.openapis.org/oas/latest.html#fixed-fields-9
		attr, err := mapSchemaToAttribute(param.Name, propBehaviorChecker([]string{}), param.Schema)
		if err != nil {
			log.Printf("[WARN] error mapping param attribute %s - %s", param.Name, err.Error())
		}
		parameterAttributes = append(parameterAttributes, *attr)
	}

	attributes := deepMergeAttributes(*requestBodyAttributes, *createResponseAttributes, *readResponseAttributes, parameterAttributes)

	resourceSchema.Attributes = *attributes
	return resourceSchema, nil
}

// getResponseParameters will retrieve the parameters from an operation or return an empty slice
func getResponseParameters(op *high.Operation) []*high.Parameter {
	if op == nil {
		return []*high.Parameter{}
	}

	return op.Parameters
}

// getRequestSchemaProxy will retrieve the schema from the request body of an operation
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func getRequestSchemaProxy(op *high.Operation) *base.SchemaProxy {
	if op == nil || op.RequestBody == nil || len(op.RequestBody.Content) == 0 {
		return nil
	}

	return getSchemaFromMediaType(op.RequestBody.Content)
}

// getResponseSchemaProxy will retrieve the schema from the response body of an operation
//   - Response codes of 200 and 201, will be prioritized, then continue to the next available 2xx code
//   - Media type will default to "application/json", then continue to the next available media type with a schema
func getResponseSchemaProxy(op *high.Operation) *base.SchemaProxy {
	if op == nil || op.Responses == nil || len(op.Responses.Codes) == 0 {
		return nil
	}

	okResponse, ok := op.Responses.Codes[oas_response_code_ok]
	if ok {
		return getSchemaFromMediaType(okResponse.Content)
	}

	createdResponse, ok := op.Responses.Codes[oas_response_code_created]
	if ok {
		return getSchemaFromMediaType(createdResponse.Content)
	}

	// Guarantee the order of processing
	codes := sortedKeys(op.Responses.Codes)
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

	return nil
}

func getSchemaFromMediaType(mediaTypes map[string]*high.MediaType) *base.SchemaProxy {
	jsonMediaType, ok := mediaTypes[oas_mediatype_json]
	if ok && jsonMediaType.Schema != nil {
		return jsonMediaType.Schema
	}

	// Guarantee the order of processing
	mediaTypeKeys := sortedKeys(mediaTypes)
	for _, key := range mediaTypeKeys {
		mediaType := mediaTypes[key]
		if mediaType.Schema != nil {
			return mediaType.Schema
		}
	}

	return nil
}

type behaviorChecker func(string, *base.Schema) ir.ComputedOptionalRequired

func propBehaviorChecker(requiredProps []string) behaviorChecker {
	requiredMap := map[string]bool{}
	for _, prop := range requiredProps {
		requiredMap[prop] = true
	}

	return func(s string, _ *base.Schema) ir.ComputedOptionalRequired {
		_, isRequired := requiredMap[s]
		if isRequired {
			return ir.Required
		} else {
			return ir.ComputedOptional
		}
	}
}

// retrieveType will return the JSON schema type. Support for multi-types is restricted to combinations of "null" and another type, i.e. ["null", "string"]
func retrieveType(schema *base.Schema) (string, error) {
	switch len(schema.Type) {
	case 0:
		return "", errors.New("property does not have a type, attribute cannot be created")
	case 1:
		return schema.Type[0], nil
	case 2:
		if schema.Type[0] == oas_type_null {
			return schema.Type[1], nil
		} else if schema.Type[1] == oas_type_null {
			return schema.Type[0], nil
		}
	}

	return "", fmt.Errorf("%w - %v", errMultiTypeSchema, schema.Type)
}

// buildSchema is a helper that builds the schema and handles nullable schemas, implemented with oneOf/anyOf OAS keywords
func buildSchema(proxy *base.SchemaProxy) (*base.Schema, error) {
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema proxy - %w", err)
	}

	if len(schema.AnyOf) == 2 {
		return getNullableSchema(schema.AnyOf[0], schema.AnyOf[1])
	}

	if len(schema.OneOf) == 2 {
		return getNullableSchema(schema.OneOf[0], schema.OneOf[1])
	}

	return schema, nil
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

	if firstType == oas_type_null {
		return secondSchema, nil
	} else if secondType == oas_type_null {
		return firstSchema, nil
	}

	return nil, fmt.Errorf("%w - %s %s", errMultiTypeSchema, firstType, secondType)
}

// Generics? ☜(ಠ_ಠ☜)
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0)

	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}
