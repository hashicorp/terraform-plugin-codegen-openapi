package datasource_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/config"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/explorer"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/datasource"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/util"
	"testing"

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
		want               []ir.DataSourceAttribute
	}{
		"merge primitives across all ops": {
			readParams: []*high.Parameter{
				{
					Name:     "string_prop",
					Required: true,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "hey this is a string, required!",
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
			want: []ir.DataSourceAttribute{
				{
					Name: "string_prop",
					String: &ir.DataSourceStringAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is a string, required!"),
						Sensitive:                pointer(true),
					},
				},
				{
					Name: "bool_prop",
					Bool: &ir.DataSourceBoolAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is a bool, required!"),
					},
				},
				{
					Name: "float64_prop",
					Float64: &ir.DataSourceFloat64Attribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey this is a float64!"),
					},
				},
				{
					Name: "number_prop",
					Number: &ir.DataSourceNumberAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey this is a number!"),
					},
				},
			},
		},
		"deep merge single nested object": {
			readParams: []*high.Parameter{
				{
					Name:     "nested_object_one",
					In:       "query",
					Required: true,
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey this is an object, required!",
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
			want: []ir.DataSourceAttribute{
				{
					Name: "nested_object_one",
					SingleNested: &ir.DataSourceSingleNestedAttribute{
						Attributes: []ir.DataSourceAttribute{
							{
								Name: "nested_object_two",
								SingleNested: &ir.DataSourceSingleNestedAttribute{
									Attributes: []ir.DataSourceAttribute{
										{
											Name: "bool_prop",
											Bool: &ir.DataSourceBoolAttribute{
												ComputedOptionalRequired: ir.ComputedOptional,
												Description:              pointer("hey this is a bool!"),
											},
										},
										{
											Name: "int64_prop",
											Int64: &ir.DataSourceInt64Attribute{
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey this is a integer!"),
											},
										},
										{
											Name: "number_prop",
											Number: &ir.DataSourceNumberAttribute{
												ComputedOptionalRequired: ir.ComputedOptional,
												Description:              pointer("hey this is a number!"),
											},
										},
									},
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey this is an object!"),
								},
							},
							{
								Name: "string_prop",
								String: &ir.DataSourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey this is a string!"),
									Sensitive:                pointer(true),
								},
							},
						},
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is an object, required!"),
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
			want: []ir.DataSourceAttribute{
				{
					Name: "array_prop",
					ListNested: &ir.DataSourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is an array, required!"),
						NestedObject: ir.DataSourceAttributeNestedObject{
							Attributes: []ir.DataSourceAttribute{
								{
									Name: "nested_array_prop",
									ListNested: &ir.DataSourceListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey this is a nested array, required!"),
										NestedObject: ir.DataSourceAttributeNestedObject{
											Attributes: []ir.DataSourceAttribute{
												{
													Name: "super_nested_bool_one",
													Bool: &ir.DataSourceBoolAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey this is a boolean!"),
													},
												},
												{
													Name: "super_nested_bool_two",
													Bool: &ir.DataSourceBoolAttribute{
														ComputedOptionalRequired: ir.Required,
														Description:              pointer("hey this is a boolean, required!"),
													},
												},
												{
													Name: "super_nested_int64",
													Int64: &ir.DataSourceInt64Attribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey this is a integer!"),
													},
												},
												{
													Name: "super_nested_string",
													String: &ir.DataSourceStringAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey this is a string!"),
														Sensitive:                pointer(false),
													},
												},
											},
										},
									},
								},
								{
									Name: "number_prop",
									Number: &ir.DataSourceNumberAttribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey this is a number!"),
									},
								},
								{
									Name: "float64_prop",
									Float64: &ir.DataSourceFloat64Attribute{
										ComputedOptionalRequired: ir.ComputedOptional,
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
			want: []ir.DataSourceAttribute{
				{
					Name: "array_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey this is an array!"),
						ElementType: ir.ElementType{
							List: &ir.ListElement{
								ElementType: &ir.ElementType{
									Object: []ir.ObjectElement{
										{
											Name: "deep_nested_list",
											ElementType: &ir.ElementType{
												List: &ir.ListElement{
													ElementType: &ir.ElementType{
														Object: []ir.ObjectElement{
															{
																Name: "deep_deep_nested_object",
																ElementType: &ir.ElementType{
																	Object: []ir.ObjectElement{
																		{
																			Name: "deep_deep_nested_bool",
																			ElementType: &ir.ElementType{
																				Bool: &ir.BoolElement{},
																			},
																		},
																		{
																			Name: "deep_deep_nested_string",
																			ElementType: &ir.ElementType{
																				String: &ir.StringElement{},
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
										{
											Name: "deep_nested_bool",
											ElementType: &ir.ElementType{
												Bool: &ir.BoolElement{},
											},
										},
										{
											Name: "deep_nested_int64",
											ElementType: &ir.ElementType{
												Int64: &ir.Int64Element{},
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
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mapper := datasource.NewDataSourceMapper(map[string]explorer.DataSource{
				"test_datasource": {
					ReadOp: createTestReadOp(testCase.readResponseSchema, testCase.readParams),
				},
			}, config.Config{})
			got, err := mapper.MapToIR()
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

func createTestReadOp(response *base.SchemaProxy, params []*high.Parameter) *high.Operation {
	return &high.Operation{
		Responses: &high.Responses{
			Codes: map[string]*high.Response{
				"200": {
					Content: map[string]*high.MediaType{
						"application/json": {
							Schema: response,
						},
					},
				},
			},
		},
		Parameters: params,
	}
}

func pointer[T any](value T) *T {
	return &value
}
