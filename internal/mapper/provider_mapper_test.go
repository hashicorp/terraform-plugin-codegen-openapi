// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: add tests for error handling/skipping bad data sources

func TestProviderMapper_basic(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		exploredProvider explorer.Provider
		want             *provider.Provider
	}{
		"provider with no schema": {
			exploredProvider: explorer.Provider{
				Name: "example",
			},
			want: &provider.Provider{
				Name: "example",
			},
		},
		"provider with schema - primitives": {
			exploredProvider: explorer.Provider{
				Name: "example",
				SchemaProxy: base.CreateSchemaProxy(&base.Schema{
					Type:     []string{"object"},
					Required: []string{"string_prop", "bool_prop"},
					Properties: map[string]*base.SchemaProxy{
						"bool_prop": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"boolean"},
							Description: "hey this is a bool, required!",
						}),
						"number_prop": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"number"},
							Description: "hey this is a number!",
						}),
						"string_prop": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"string"},
							Format:      util.OAS_format_password,
							Description: "hey this is a string, required!",
						}),
						"float64_prop": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"number"},
							Format:      "float",
							Description: "hey this is a float64!",
						}),
					},
				}),
			},
			want: &provider.Provider{
				Name: "example",
				Schema: &provider.Schema{
					Attributes: provider.Attributes{
						{
							Name: "bool_prop",
							Bool: &provider.BoolAttribute{
								OptionalRequired: schema.Required,
								Description:      pointer("hey this is a bool, required!"),
							},
						},
						{
							Name: "float64_prop",
							Float64: &provider.Float64Attribute{
								OptionalRequired: schema.Optional,
								Description:      pointer("hey this is a float64!"),
							},
						},
						{
							Name: "number_prop",
							Number: &provider.NumberAttribute{
								OptionalRequired: schema.Optional,
								Description:      pointer("hey this is a number!"),
							},
						},
						{
							Name: "string_prop",
							String: &provider.StringAttribute{
								OptionalRequired: schema.Required,
								Description:      pointer("hey this is a string, required!"),
								Sensitive:        pointer(true),
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

			mapper := mapper.NewProviderMapper(testCase.exploredProvider, config.Config{})
			got, err := mapper.MapToIR()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(got, testCase.want); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
