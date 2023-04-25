package mapper

import (
	"errors"
	"fmt"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"log"
	"sort"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

var _ ResourceMapper = resourceMapper{}

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

// Merge the following to create the schema, priority as follows
// 1. Create op - request body
//
// TODO: 2. Read op - request body
//   - Any values here not already from POST are computed
//
// TODO: 3. Read op - parameters (path, query, headers?)
func generateResourceSchema(resource explorer.Resource) (*ir.ResourceSchema, error) {
	resourceSchema := &ir.ResourceSchema{
		Attributes: []ir.ResourceAttribute{},
	}

	createRequestSchema := getRequestSchemaProxy(resource.CreateOp.RequestBody)
	if createRequestSchema == nil {
		return nil, errors.New("no request schema found!")
	}

	// TODO: this assumes the create request is always an object, which is probably okay for resources
	attributes, err := mapSchemaToObjectAttributes(createRequestSchema)
	if err != nil {
		return nil, err
	}

	resourceSchema.Attributes = *attributes
	return resourceSchema, nil
}

func getRequestSchemaProxy(request *high.RequestBody) *base.SchemaProxy {
	if request.Content["application/json"] == nil {
		return nil
	}
	return request.Content["application/json"].Schema
}

type behaviorChecker func(string) ir.ComputedOptionalRequired

func propBehaviorChecker(requiredProps []string) behaviorChecker {
	requiredMap := map[string]bool{}
	for _, prop := range requiredProps {
		requiredMap[prop] = true
	}

	return func(s string) ir.ComputedOptionalRequired {
		_, isRequired := requiredMap[s]
		if isRequired {
			return ir.Required
		} else {
			return ir.Optional
		}
	}
}

func retrieveType(typeArr []string) (string, error) {
	if len(typeArr) < 1 {
		return "", fmt.Errorf("property does not have a type, therefore an attribute cannot be created")
	}
	if len(typeArr) > 1 {
		return "", fmt.Errorf("property has multiple types, %v, therefore an attribute cannot be created", typeArr)
	}

	return typeArr[0], nil
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
