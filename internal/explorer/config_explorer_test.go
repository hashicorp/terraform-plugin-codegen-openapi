// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func Test_ConfigExplorer_FindResources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathItems map[string]*high.PathItem
		config    config.Config
		want      map[string]explorer.Resource
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
			pathItems: map[string]*high.PathItem{
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
			},
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
			pathItems: map[string]*high.PathItem{
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
			},
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
		"non-existent paths and methods are ignored ": {
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
							Path:   "/fakepath",
							Method: "PUT",
						},
						Delete: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "FAKEMETHOD",
						},
					},
				},
			},
			pathItems: map[string]*high.PathItem{
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
			},
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
						AttributeOptions: explorer.AttributeOptions{
							Overrides: map[string]explorer.Override{},
						},
					},
				},
			},
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
			pathItems: map[string]*high.PathItem{
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
			},
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
			got, err := explorer.FindResources()

			if err != nil {
				t.Fatalf("was not expecting error, got: %s", err)
			}

			if diff := cmp.Diff(got, testCase.want, cmpopts.IgnoreUnexported(high.Operation{})); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func Test_ConfigExplorer_FindDataSources(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pathItems map[string]*high.PathItem
		config    config.Config
		want      map[string]explorer.DataSource
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
			pathItems: map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			},
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
			pathItems: map[string]*high.PathItem{
				"/resources/two/{resource_id}": {
					Head: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			},
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
		"non-existent paths and methods are ignored ": {
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
			pathItems: map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			},
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
		"schema options pass-through": {
			config: config.Config{
				DataSources: map[string]config.DataSource{
					"test_resource": {
						Read: &config.OpenApiSpecLocation{
							Path:   "/resources/{resource_id}",
							Method: "GET",
						},
						SchemaOptions: config.SchemaOptions{
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
			pathItems: map[string]*high.PathItem{
				"/resources/{resource_id}": {
					Get: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
				},
			},
			want: map[string]explorer.DataSource{
				"test_resource": {
					ReadOp: &high.Operation{
						Description: "read op here",
						OperationId: "read_resource",
					},
					SchemaOptions: explorer.SchemaOptions{
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

			if err != nil {
				t.Fatalf("was not expecting error, got: %s", err)
			}

			if diff := cmp.Diff(got, testCase.want, cmpopts.IgnoreUnexported(high.Operation{})); diff != "" {
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
				Extensions:  map[string]any{},
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

				if diff := cmp.Diff(gotSchema, testCase.expectedSchema, cmpopts.IgnoreUnexported(base.Schema{})); diff != "" {
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
