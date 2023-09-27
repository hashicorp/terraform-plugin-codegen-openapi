// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_test

import (
	"log/slog"
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// TODO: add tests for error handling/skipping bad data sources

func TestDataSourceMapper_basic_merges(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		readResponseSchema *base.SchemaProxy
		readParams         []*high.Parameter
		schemaOptions      explorer.SchemaOptions
		want               datasource.Attributes
	}{
		"merge primitives across all ops": {
			readParams: []*high.Parameter{
				{
					Name:        "string_prop",
					Required:    true,
					In:          "path",
					Description: "hey this is a string, required and overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
				{
					Name:     "bool_prop",
					Required: true,
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
				},
				{
					Name: "float64_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey this is a float64!",
					}),
				},
			},
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool!",
					}),
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey this is a number!",
					}),
				},
			}),
			want: datasource.Attributes{
				{
					Name: "string_prop",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is a string, required and overidden!"),
						Sensitive:                pointer(true),
					},
				},
				{
					Name: "bool_prop",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is a bool, required!"),
					},
				},
				{
					Name: "float64_prop",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is a float64!"),
					},
				},
				{
					Name: "number_prop",
					Number: &datasource.NumberAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("hey this is a number!"),
					},
				},
			},
		},
		"deep merge single nested object": {
			readParams: []*high.Parameter{
				{
					Name:        "nested_object_one",
					In:          "query",
					Required:    true,
					Description: "hey this is an object, required + overidden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "you shouldn't see this because the description is overridden!",
						Properties: map[string]*base.SchemaProxy{
							"nested_object_two": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"int64_prop"},
								Description: "hey this is an object!",
								Properties: map[string]*base.SchemaProxy{
									"bool_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"boolean"},
										Description: "hey this is a bool!",
									}),
									"int64_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"integer"},
										Description: "hey this is a integer!",
									}),
								},
							}),
						},
					}),
				},
			},
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_object_one": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "this one already exists, so you shouldn't see this description!",
						Properties: map[string]*base.SchemaProxy{
							"string_prop": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Format:      util.OAS_format_password,
								Description: "hey this is a string!",
							}),
							"nested_object_two": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Description: "this one already exists, so you shouldn't see this description!",
								Properties: map[string]*base.SchemaProxy{
									"bool_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"boolean"},
										Description: "hey this is a bool!",
									}),
									"number_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"number"},
										Description: "hey this is a number!",
									}),
								},
							}),
						},
					}),
				},
			}),
			want: datasource.Attributes{
				{
					Name: "nested_object_one",
					SingleNested: &datasource.SingleNestedAttribute{
						Attributes: []datasource.Attribute{
							{
								Name: "nested_object_two",
								SingleNested: &datasource.SingleNestedAttribute{
									Attributes: []datasource.Attribute{
										{
											Name: "bool_prop",
											Bool: &datasource.BoolAttribute{
												ComputedOptionalRequired: schema.ComputedOptional,
												Description:              pointer("hey this is a bool!"),
											},
										},
										{
											Name: "int64_prop",
											Int64: &datasource.Int64Attribute{
												ComputedOptionalRequired: schema.Required,
												Description:              pointer("hey this is a integer!"),
											},
										},
										{
											Name: "number_prop",
											Number: &datasource.NumberAttribute{
												ComputedOptionalRequired: schema.Computed,
												Description:              pointer("hey this is a number!"),
											},
										},
									},
									ComputedOptionalRequired: schema.ComputedOptional,
									Description:              pointer("hey this is an object!"),
								},
							},
							{
								Name: "string_prop",
								String: &datasource.StringAttribute{
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("hey this is a string!"),
									Sensitive:                pointer(true),
								},
							},
						},
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is an object, required + overidden!"),
					},
				},
			},
		},
		"deep merge list nested array": {
			readParams: []*high.Parameter{
				{
					Name:     "array_prop",
					In:       "query",
					Required: true,
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey this is an array, required!",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_array_prop"},
								Properties: map[string]*base.SchemaProxy{
									"nested_array_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"array"},
										Description: "hey this is a nested array, required!",
										Items: &base.DynamicValue[*base.SchemaProxy, bool]{
											A: base.CreateSchemaProxy(&base.Schema{
												Type:     []string{"object"},
												Required: []string{"super_nested_bool_two"},
												Properties: map[string]*base.SchemaProxy{
													"super_nested_bool_one": base.CreateSchemaProxy(&base.Schema{
														Type:        []string{"boolean"},
														Description: "hey this is a boolean!",
													}),
													"super_nested_bool_two": base.CreateSchemaProxy(&base.Schema{
														Type:        []string{"boolean"},
														Description: "hey this is a boolean, required!",
													}),
													"super_nested_int64": base.CreateSchemaProxy(&base.Schema{
														Type:        []string{"integer"},
														Description: "hey this is a integer!",
													}),
												},
											}),
										},
									}),
									"number_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"number"},
										Description: "hey this is a number!",
									}),
								},
							}),
						},
					}),
				},
			},
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"array_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey this is an array!",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_array_prop"},
								Properties: map[string]*base.SchemaProxy{
									"float64_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"number"},
										Format:      "double",
										Description: "hey this is a float64!",
									}),
									"nested_array_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"array"},
										Description: "hey this is a nested array, required!",
										Items: &base.DynamicValue[*base.SchemaProxy, bool]{
											A: base.CreateSchemaProxy(&base.Schema{
												Type: []string{"object"},
												Properties: map[string]*base.SchemaProxy{
													"super_nested_string": base.CreateSchemaProxy(&base.Schema{
														Type:        []string{"string"},
														Description: "hey this is a string!",
													}),
												},
											}),
										},
									}),
								},
							}),
						},
					}),
				},
			}),
			want: datasource.Attributes{
				{
					Name: "array_prop",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is an array, required!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_array_prop",
									ListNested: &datasource.ListNestedAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey this is a nested array, required!"),
										NestedObject: datasource.NestedAttributeObject{
											Attributes: []datasource.Attribute{
												{
													Name: "super_nested_bool_one",
													Bool: &datasource.BoolAttribute{
														ComputedOptionalRequired: schema.ComputedOptional,
														Description:              pointer("hey this is a boolean!"),
													},
												},
												{
													Name: "super_nested_bool_two",
													Bool: &datasource.BoolAttribute{
														ComputedOptionalRequired: schema.Required,
														Description:              pointer("hey this is a boolean, required!"),
													},
												},
												{
													Name: "super_nested_int64",
													Int64: &datasource.Int64Attribute{
														ComputedOptionalRequired: schema.ComputedOptional,
														Description:              pointer("hey this is a integer!"),
													},
												},
												{
													Name: "super_nested_string",
													String: &datasource.StringAttribute{
														ComputedOptionalRequired: schema.Computed,
														Description:              pointer("hey this is a string!"),
													},
												},
											},
										},
									},
								},
								{
									Name: "number_prop",
									Number: &datasource.NumberAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey this is a number!"),
									},
								},
								{
									Name: "float64_prop",
									Float64: &datasource.Float64Attribute{
										ComputedOptionalRequired: schema.Computed,
										Description:              pointer("hey this is a float64!"),
									},
								},
							},
						},
					},
				},
			},
		},
		"deep merge list array with object element types": {
			readParams: []*high.Parameter{
				{
					Name: "array_prop",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey this is an array!",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: map[string]*base.SchemaProxy{
											"deep_nested_list": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"array"},
												Items: &base.DynamicValue[*base.SchemaProxy, bool]{
													A: base.CreateSchemaProxy(&base.Schema{
														Type: []string{"object"},
														Properties: map[string]*base.SchemaProxy{
															"deep_deep_nested_object": base.CreateSchemaProxy(&base.Schema{
																Type: []string{"object"},
																Properties: map[string]*base.SchemaProxy{
																	"deep_deep_nested_bool": base.CreateSchemaProxy(&base.Schema{
																		Type: []string{"boolean"},
																	}),
																},
															}),
														},
													}),
												},
											}),
										},
									}),
								},
							}),
						},
					}),
				},
			},
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"array_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey this is an array!",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: map[string]*base.SchemaProxy{
											"deep_nested_list": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"array"},
												Items: &base.DynamicValue[*base.SchemaProxy, bool]{
													A: base.CreateSchemaProxy(&base.Schema{
														Type: []string{"object"},
														Properties: map[string]*base.SchemaProxy{
															"deep_deep_nested_object": base.CreateSchemaProxy(&base.Schema{
																Type: []string{"object"},
																Properties: map[string]*base.SchemaProxy{
																	"deep_deep_nested_string": base.CreateSchemaProxy(&base.Schema{
																		Type: []string{"string"},
																	}),
																},
															}),
														},
													}),
												},
											}),
											"deep_nested_bool": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"deep_nested_int64": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"integer"},
											}),
										},
									}),
								},
							}),
						},
					}),
				},
			}),
			want: datasource.Attributes{
				{
					Name: "array_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is an array!"),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "deep_nested_list",
												List: &schema.ListType{
													ElementType: schema.ElementType{
														Object: &schema.ObjectType{
															AttributeTypes: []schema.ObjectAttributeType{
																{
																	Name: "deep_deep_nested_object",
																	Object: &schema.ObjectType{
																		AttributeTypes: []schema.ObjectAttributeType{
																			{
																				Name: "deep_deep_nested_bool",
																				Bool: &schema.BoolType{},
																			},
																			{
																				Name:   "deep_deep_nested_string",
																				String: &schema.StringType{},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											{
												Name: "deep_nested_bool",
												Bool: &schema.BoolType{},
											},
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
			},
		},
		"precedence and configurability": {
			readParams: []*high.Parameter{
				{
					Name:     "read_parameter_optional_read_parameter_only",
					Required: false,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "read_parameter_optional_read_response",
					Required: false,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "read_parameter_required_read_parameter_only",
					Required: true,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "read_parameter_required_read_response",
					Required: true,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			},
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					// Simulate API returning parameter in response
					"read_parameter_optional_read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					// Simulate API returning parameter in response
					"read_parameter_required_read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
			want: datasource.Attributes{
				{
					Name: "read_parameter_optional_read_parameter_only",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "read_parameter_optional_read_response",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "read_parameter_required_read_parameter_only",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "read_parameter_required_read_response",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "read_response",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Computed,
					},
				},
			},
		},
		"parameter match for path and query params": {
			readParams: []*high.Parameter{
				{
					Name:     "read_path_parameter",
					Required: true,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "read_query_parameter",
					Required: false,
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"boolean"},
					}),
				},
			},
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"attribute_required": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"attribute_computed_optional": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"boolean"},
					}),
				},
			}),
			schemaOptions: explorer.SchemaOptions{
				AttributeOptions: explorer.AttributeOptions{
					Aliases: map[string]string{
						"read_path_parameter":  "attribute_required",
						"read_query_parameter": "attribute_computed_optional",
					},
				},
			},
			want: datasource.Attributes{
				{
					Name: "attribute_required",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "attribute_computed_optional",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mapper := mapper.NewDataSourceMapper(map[string]explorer.DataSource{
				"test_datasource": {
					ReadOp:        createTestReadOp(testCase.readResponseSchema, testCase.readParams),
					SchemaOptions: testCase.schemaOptions,
				},
			}, config.Config{})
			got, err := mapper.MapToIR(slog.Default())
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(got) != 1 {
				t.Fatalf("expected only one DataSource, got: %d", len(got))
			}

			if diff := cmp.Diff(got[0].Schema.Attributes, testCase.want); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
