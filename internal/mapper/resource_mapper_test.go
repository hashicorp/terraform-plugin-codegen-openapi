package mapper_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// TODO: add tests for request / response combining + multiple resources
// TODO: add tests for error handling/skipping bad resources

func createTestResources(requestSchema *base.Schema) map[string]explorer.Resource {
	requestSchemaProxy := base.CreateSchemaProxy(requestSchema)

	return map[string]explorer.Resource{
		"test_resource": {
			CreateOp: &high.Operation{
				RequestBody: &high.RequestBody{
					Content: map[string]*high.MediaType{
						"application/json": {
							Schema: requestSchemaProxy,
						},
					},
				},
			},
		},
	}
}

func pointer[T any](value T) *T {
	return &value
}
