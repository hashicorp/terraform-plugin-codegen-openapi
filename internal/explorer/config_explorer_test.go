// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"gopkg.in/yaml.v3"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func Test_ConfigExplorer_FindResources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathItems   *orderedmap.Map[string, *high.PathItem]
		config      config.Config
		want        map[string]explorer.Resource
		expectedErr error
	}{
		"valid CRUD ops": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Create: &config.OpenApiSpecLocation{
							Path:   "/resources",
							Method: "POST",
						},
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "GET",
						},
						Update: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "PUT",
						},
						Delete: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "DELETE",
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources": {
					Post: &high.Operation{
						Description: "create op here",
						OperationId: "create_resource",
					},
				},
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					Put: &high.Operation{
						Description: "update op here",
						OperationId: "update_resource",
					},
					Delete: &high.Operation{
						Description: "delete op here",
						OperationId: "delete_resource",
					},
				},
			}),
			want: map[string]explorer.Resource{
				"test_resource": {
					CreateOp: &high.Operation{
						Description: "create op here",
						OperationId: "create_resource",
					},
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					UpdateOp: &high.Operation{
						Description: "update op here",
						OperationId: "update_resource",
					},
					DeleteOp: &high.Operation{
						Description: "delete op here",
						OperationId: "delete_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
						AttributeOptions: explorer.AttributeOptions{
							Overrides: map[string]explorer.Override{},
						},
					},
				},
			},
		},
		"valid alternative CRUD ops - options, head, patch, trace": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Create: &config.OpenApiSpecLocation{
							Path:   "/resources/one",
							Method: "OPTIONS",
						},
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/two/{resource_id}",
							Method: "HEAD",
						},
						Update: &config.OpenApiSpecLocation{
							Path:   "/resources/three/{resource_id}",
							Method: "PATCH",
						},
						Delete: &config.OpenApiSpecLocation{
							Path:   "/resources/one",
							Method: "TRACE",
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources/one": {
					Options: &high.Operation{
						Description: "create op here",
						OperationId: "create_resource",
					},
					Trace: &high.Operation{
						Description: "delete op here",
						OperationId: "delete_resource",
					},
				},
				"/resources/two/{resource_id}": {
					Head: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
				"/resources/three/{resource_id}": {
					Patch: &high.Operation{
						Description: "update op here",
						OperationId: "update_resource",
					},
				},
			}),
			want: map[string]explorer.Resource{
				"test_resource": {
					CreateOp: &high.Operation{
						Description: "create op here",
						OperationId: "create_resource",
					},
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					UpdateOp: &high.Operation{
						Description: "update op here",
						OperationId: "update_resource",
					},
					DeleteOp: &high.Operation{
						Description: "delete op here",
						OperationId: "delete_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
						AttributeOptions: explorer.AttributeOptions{
							Overrides: map[string]explorer.Override{},
						},
					},
				},
			},
		},
		"non-existent create path throws error": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Create: &config.OpenApiSpecLocation{
							Path:   "/fakepath",
							Method: "POST",
						},
					},
				},
			},
			pathItems:   orderedmap.ToOrderedMap(map[string]*high.PathItem{}),
			expectedErr: errors.New(`failed to extract 'test_resource.create': path '/fakepath' not found in OpenAPI spec`),
		},
		"non-existent read path throws error": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/fakepath",
							Method: "GET",
						},
					},
				},
			},
			pathItems:   orderedmap.ToOrderedMap(map[string]*high.PathItem{}),
			expectedErr: errors.New(`failed to extract 'test_resource.read': path '/fakepath' not found in OpenAPI spec`),
		},
		"non-existent update path throws error": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Update: &config.OpenApiSpecLocation{
							Path:   "/fakepath",
							Method: "PUT",
						},
					},
				},
			},
			pathItems:   orderedmap.ToOrderedMap(map[string]*high.PathItem{}),
			expectedErr: errors.New(`failed to extract 'test_resource.update': path '/fakepath' not found in OpenAPI spec`),
		},
		"non-existent delete path throws error": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Delete: &config.OpenApiSpecLocation{
							Path:   "/fakepath",
							Method: "DELETE",
						},
					},
				},
			},
			pathItems:   orderedmap.ToOrderedMap(map[string]*high.PathItem{}),
			expectedErr: errors.New(`failed to extract 'test_resource.delete': path '/fakepath' not found in OpenAPI spec`),
		},
		"non-existent method throws error": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Update: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "FAKE",
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Put: &high.Operation{
						Description: "update op here",
						OperationId: "update_resource",
					},
				},
			}),
			expectedErr: errors.New(`failed to extract 'test_resource.update': method 'FAKE' not found at OpenAPI path '/resources/{resource_id}'`),
		},
		"schema options pass-through": {
			config: config.Config{
				Resources: map[string]config.Resource{
					"test_resource": {
						Create: &config.OpenApiSpecLocation{
							Path:   "/resources",
							Method: "POST",
						},
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "GET",
						},
						SchemaOptions: config.SchemaOptions{
							Ignores: []string{"ignore1.abc", "ignore2.def"},
							AttributeOptions: config.AttributeOptions{
								Aliases: map[string]string{
									"otherId": "id",
								},
								Overrides: map[string]config.Override{
									"test": {
										Description:              "test description for override",
										ComputedOptionalRequired: "computed_optional",
									},
								},
							},
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources": {
					Post: &high.Operation{
						Description: "create op here",
						OperationId: "create_resource",
					},
				},
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			}),
			want: map[string]explorer.Resource{
				"test_resource": {
					CreateOp: &high.Operation{
						Description: "create op here",
						OperationId: "create_resource",
					},
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
						Ignores: []string{"ignore1.abc", "ignore2.def"},
						AttributeOptions: explorer.AttributeOptions{
							Aliases: map[string]string{
								"otherId": "id",
							},
							Overrides: map[string]explorer.Override{
								"test": {
									Description:              "test description for override",
									ComputedOptionalRequired: "computed_optional",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			explorer := explorer.NewConfigExplorer(high.Document{Paths: &high.Paths{PathItems: testCase.pathItems}}, testCase.config)
			got, err := explorer.FindResources()

			if testCase.expectedErr != nil {
				if err == nil {
					t.Fatal("expected an error, but got none")
				}

				if testCase.expectedErr.Error() != err.Error() {
					t.Fatalf("expected err: %s, got: %s", testCase.expectedErr, err)
				}

				return
			}

			if diff := cmp.Diff(got, testCase.want, cmpopts.IgnoreUnexported(high.Operation{})); testCase.expectedErr == nil && diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func Test_ConfigExplorer_FindDataSources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathItems   *orderedmap.Map[string, *high.PathItem]
		config      config.Config
		want        map[string]explorer.DataSource
		expectedErr error
	}{
		"valid read op": {
			config: config.Config{
				DataSources: map[string]config.DataSource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "GET",
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			}),
			want: map[string]explorer.DataSource{
				"test_resource": {
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
						AttributeOptions: explorer.AttributeOptions{
							Overrides: map[string]explorer.Override{},
						},
					},
				},
			},
		},
		"valid read op - alternative methods": {
			config: config.Config{
				DataSources: map[string]config.DataSource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/two/{resource_id}",
							Method: "HEAD",
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources/two/{resource_id}": {
					Head: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			}),
			want: map[string]explorer.DataSource{
				"test_resource": {
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
						AttributeOptions: explorer.AttributeOptions{
							Overrides: map[string]explorer.Override{},
						},
					},
				},
			},
		},
		"non-existent read path throws error": {
			config: config.Config{
				DataSources: map[string]config.DataSource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/fakepath",
							Method: "GET",
						},
					},
				},
			},
			pathItems:   orderedmap.ToOrderedMap(map[string]*high.PathItem{}),
			expectedErr: errors.New(`failed to extract 'test_resource.read': path '/fakepath' not found in OpenAPI spec`),
		},
		"non-existent method throws error": {
			config: config.Config{
				DataSources: map[string]config.DataSource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "FAKE",
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			}),
			expectedErr: errors.New(`failed to extract 'test_resource.read': method 'FAKE' not found at OpenAPI path '/resources/{resource_id}'`),
		},
		"schema options pass-through": {
			config: config.Config{
				DataSources: map[string]config.DataSource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "GET",
						},
						SchemaOptions: config.SchemaOptions{
							Ignores: []string{"ignore1.abc", "ignore2.def"},
							AttributeOptions: config.AttributeOptions{
								Aliases: map[string]string{
									"otherId": "id",
								},
								Overrides: map[string]config.Override{
									"test": {
										Description: "test description for override",
									},
								},
							},
						},
					},
				},
			},
			pathItems: orderedmap.ToOrderedMap(map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			}),
			want: map[string]explorer.DataSource{
				"test_resource": {
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
						Ignores: []string{"ignore1.abc", "ignore2.def"},
						AttributeOptions: explorer.AttributeOptions{
							Aliases: map[string]string{
								"otherId": "id",
							},
							Overrides: map[string]explorer.Override{
								"test": {
									Description: "test description for override",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			explorer := explorer.NewConfigExplorer(high.Document{Paths: &high.Paths{PathItems: testCase.pathItems}}, testCase.config)
			got, err := explorer.FindDataSources()

			if testCase.expectedErr != nil {
				if err == nil {
					t.Fatal("expected an error, but got none")
				}

				if testCase.expectedErr.Error() != err.Error() {
					t.Fatalf("expected err: %s, got: %s", testCase.expectedErr, err)
				}

				return
			}

			if diff := cmp.Diff(got, testCase.want, cmpopts.IgnoreUnexported(high.Operation{})); testCase.expectedErr == nil && diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func Test_ConfigExplorer_FindProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config         config.Config
		expectedName   string
		expectedSchema *base.Schema
	}{
		"valid provider name from config": {
			config: config.Config{
				Provider: config.Provider{
					Name: "heres_the_provider_name",
				},
			},
			expectedName: "heres_the_provider_name",
		},
		"valid and resolvable schema_ref from config": {
			config: config.Config{
				Provider: config.Provider{
					Name:      "example",
					SchemaRef: "#/components/schemas/example_provider",
				},
			},
			expectedName: "example",
			// We only really care that it resolves the right schema for this logic, so just comparing description/type
			expectedSchema: &base.Schema{
				Type:        []string{"object"},
				Description: "This is the provider schema",
				Extensions:  orderedmap.ToOrderedMap(map[string]*yaml.Node{}),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			oasModel, err := buildTestOAS()
			if err != nil {
				t.Fatal(err)
			}

			explorer := explorer.NewConfigExplorer(oasModel, testCase.config)
			got, err := explorer.FindProvider()

			if err != nil {
				t.Fatalf("was not expecting error, got: %s", err)
			}

			if got.Name != testCase.expectedName {
				t.Fatalf("expected provider name %s, got: %s", testCase.expectedName, got.Name)
			}

			if testCase.expectedSchema == nil && got.SchemaProxy != nil {
				t.Fatal("expected schema proxy to be empty")
			}

			if testCase.expectedSchema != nil {
				if got.SchemaProxy == nil {
					t.Fatal("expected a schema proxy for provider, but didn't return one")
				}

				gotSchema, err := got.SchemaProxy.BuildSchema()
				if err != nil {
					t.Fatalf("error building returned schema proxy: %s", err)
				}

				if diff := cmp.Diff(gotSchema, testCase.expectedSchema, cmpopts.IgnoreUnexported(base.Schema{}), cmpopts.IgnoreFields(base.Schema{}, "Extensions")); diff != "" {
					t.Errorf("unexpected difference: %s", diff)
				}
			}
		})
	}
}

func buildTestOAS() (high.Document, error) {
	testOAS := `
openapi: 3.1.0
info:
  title: Example API
components:
  schemas:
    example_provider:
      description: This is the provider schema
      type: object
    not_the_provider:
      description: Not this one
      type: object
`

	doc, err := libopenapi.NewDocument([]byte(testOAS))
	if err != nil {
		return high.Document{}, fmt.Errorf("unexpected error parsing test OAS: %w", err)
	}

	testOASModel, errs := doc.BuildV3Model()
	if len(errs) > 0 {
		var errResult error
		for _, err := range errs {
			errResult = errors.Join(errResult, err)
		}
		return high.Document{}, fmt.Errorf("unexpected error building test OAS: %w", errResult)
	}

	return testOASModel.Model, nil
}
