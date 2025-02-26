// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_test

import (
	"log/slog"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
)

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
					Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
					}),
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
		"provider with schema - ignores": {
			exploredProvider: explorer.Provider{
				Name: "example",
				Ignores: []string{
					"bool_prop",
					"nested_obj.bool_prop",
					"nested_array.deep_nested_bool",
					"nested_map.deep_nested_bool",
				},
				SchemaProxy: base.CreateSchemaProxy(&base.Schema{
					Type: []string{"object"},
					Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
						"bool_prop": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"boolean"},
							Description: "This boolean is going to be ignored!",
						}),
						"number_prop": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"number"},
							Description: "hey this is a number!",
						}),
						"nested_obj": base.CreateSchemaProxy(&base.Schema{
							Type: []string{"object"},
							Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
								"bool_prop": base.CreateSchemaProxy(&base.Schema{
									Type:        []string{"boolean"},
									Description: "This boolean is going to be ignored!",
								}),
								"string_prop": base.CreateSchemaProxy(&base.Schema{
									Type:        []string{"string"},
									Description: "hey this is a string!",
								}),
							}),
						}),
						"nested_array": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"array"},
							Description: "hey this is an array!",
							Items: &base.DynamicValue[*base.SchemaProxy, bool]{
								A: base.CreateSchemaProxy(&base.Schema{
									Type: []string{"array"},
									Items: &base.DynamicValue[*base.SchemaProxy, bool]{
										A: base.CreateSchemaProxy(&base.Schema{
											Type: []string{"object"},
											Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
												"deep_nested_bool": base.CreateSchemaProxy(&base.Schema{
													Type: []string{"boolean"},
												}),
												"deep_nested_int64": base.CreateSchemaProxy(&base.Schema{
													Type: []string{"integer"},
												}),
											}),
										}),
									},
								}),
							},
						}),
						"nested_map": base.CreateSchemaProxy(&base.Schema{
							Type:        []string{"object"},
							Description: "hey this is a map!",
							AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
								A: base.CreateSchemaProxy(&base.Schema{
									Type: []string{"object"},
									Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
										"deep_nested_bool": base.CreateSchemaProxy(&base.Schema{
											Type: []string{"boolean"},
										}),
										"deep_nested_int64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"integer"},
											Description: "hey this is an int64!",
										}),
									}),
								}),
							},
						}),
					}),
				}),
			},
			want: &provider.Provider{
				Name: "example",
				Schema: &provider.Schema{
					Attributes: provider.Attributes{
						{
							Name: "nested_array",
							List: &provider.ListAttribute{
								OptionalRequired: schema.Optional,
								Description:      pointer("hey this is an array!"),
								ElementType: schema.ElementType{
									List: &schema.ListType{
										ElementType: schema.ElementType{
											Object: &schema.ObjectType{
												AttributeTypes: []schema.ObjectAttributeType{
													{
														Name:  "deep_nested_int64",
														Int64: &schema.Int64Type{},
													},
												},
											},
										},
									},
								},
							},
						},
						{
							Name: "nested_map",
							MapNested: &provider.MapNestedAttribute{
								OptionalRequired: schema.Optional,
								Description:      pointer("hey this is a map!"),
								NestedObject: provider.NestedAttributeObject{
									Attributes: []provider.Attribute{
										{
											Name: "deep_nested_int64",
											Int64: &provider.Int64Attribute{
												OptionalRequired: schema.Optional,
												Description:      pointer("hey this is an int64!"),
											},
										},
									},
								},
							},
						},
						{
							Name: "nested_obj",
							SingleNested: &provider.SingleNestedAttribute{
								OptionalRequired: schema.Optional,
								Attributes: []provider.Attribute{
									{
										Name: "string_prop",
										String: &provider.StringAttribute{
											OptionalRequired: schema.Optional,
											Description:      pointer("hey this is a string!"),
										},
									},
								},
							},
						},
						{
							Name: "number_prop",
							Number: &provider.NumberAttribute{
								OptionalRequired: schema.Optional,
								Description:      pointer("hey this is a number!"),
							},
						},
					},
				},
			},
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mapper := mapper.NewProviderMapper(testCase.exploredProvider, config.Config{})
			got, err := mapper.MapToIR(slog.Default())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(got, testCase.want); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
