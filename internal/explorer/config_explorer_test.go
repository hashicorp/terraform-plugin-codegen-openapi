package explorer_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/config"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

			// Unexported high.Operation.low is throwing errors from cmp
			// TODO: this is hacky + not recommended, should see if there is a better comparison method long-term
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

			// Unexported high.Operation.low is throwing errors from cmp
			// TODO: this is hacky + not recommended, should see if there is a better comparison method long-term
			if diff := cmp.Diff(got, testCase.want, cmpopts.IgnoreUnexported(high.Operation{})); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func Test_ConfigExplorer_FindProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		config config.Config
		want   explorer.Provider
	}{
		"valid provider name from config": {
			config: config.Config{
				Provider: config.Provider{
					Name: "heres_the_provider_name",
				},
			},
			want: explorer.Provider{
				Name: "heres_the_provider_name",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			explorer := explorer.NewConfigExplorer(high.Document{}, testCase.config)
			got, err := explorer.FindProvider()

			if err != nil {
				t.Fatalf("was not expecting error, got: %s", err)
			}

			if diff := cmp.Diff(got, testCase.want); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
