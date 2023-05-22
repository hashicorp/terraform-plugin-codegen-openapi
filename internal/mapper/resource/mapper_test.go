package resource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/resource"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// TODO: add tests for error handling/skipping bad resources

func TestResourceMapper_basic_merges(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		createRequestSchema  *base.SchemaProxy
		createResponseSchema *base.SchemaProxy
		readResponseSchema   *base.SchemaProxy
		readParams           []*high.Parameter
		want                 []ir.ResourceAttribute
	}{
		"merge primitives across all ops": {
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_prop", "int64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey this is a bool, required!",
					}),
					"int64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Description: "hey this is an int64, required!",
					}),
				},
			}),
			createResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"int64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Description: "this one already exists, so you shouldn't see this description!",
					}),
				},
			}),
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey this is a number!",
					}),
				},
			}),
			readParams: []*high.Parameter{
				{
					Name: "string_prop",
					In:   "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "hey this is a string!",
					}),
				},
			},
			want: []ir.ResourceAttribute{
				{
					Name: "bool_prop",
					Bool: &ir.ResourceBoolAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is a bool, required!"),
					},
				},
				{
					Name: "int64_prop",
					Int64: &ir.ResourceInt64Attribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is an int64, required!"),
					},
				},
				{
					Name: "number_prop",
					Number: &ir.ResourceNumberAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey this is a number!"),
					},
				},
				{
					Name: "string_prop",
					String: &ir.ResourceStringAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey this is a string!"),
						Sensitive:                pointer(true),
					},
				},
			},
		},
		"deep merge single nested object": {
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_object_one"},
				Properties: map[string]*base.SchemaProxy{
					"nested_object_one": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"bool_prop"},
						Description: "hey this is an object, required!",
						Properties: map[string]*base.SchemaProxy{
							"bool_prop": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"boolean"},
								Description: "hey this is a bool, required!",
							}),
						},
					}),
				},
			}),
			createResponseSchema: base.CreateSchemaProxy(&base.Schema{
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
								Description: "hey this is an object!",
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
			readParams: []*high.Parameter{
				{
					Name: "nested_object_one",
					In:   "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "this one already exists, so you shouldn't see this description!",
						Properties: map[string]*base.SchemaProxy{
							"nested_object_two": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"int64_prop"},
								Description: "this one already exists, so you shouldn't see this description!",
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
			want: []ir.ResourceAttribute{
				{
					Name: "nested_object_one",
					SingleNested: &ir.ResourceSingleNestedAttribute{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "bool_prop",
								Bool: &ir.ResourceBoolAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey this is a bool, required!"),
								},
							},
							{
								Name: "nested_object_two",
								SingleNested: &ir.ResourceSingleNestedAttribute{
									Attributes: []ir.ResourceAttribute{
										{
											Name: "bool_prop",
											Bool: &ir.ResourceBoolAttribute{
												ComputedOptionalRequired: ir.ComputedOptional,
												Description:              pointer("hey this is a bool!"),
											},
										},
										{
											Name: "number_prop",
											Number: &ir.ResourceNumberAttribute{
												ComputedOptionalRequired: ir.ComputedOptional,
												Description:              pointer("hey this is a number!"),
											},
										},
										{
											Name: "int64_prop",
											Int64: &ir.ResourceInt64Attribute{
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey this is a integer!"),
											},
										},
									},
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey this is an object!"),
								},
							},
							{
								Name: "string_prop",
								String: &ir.ResourceStringAttribute{
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
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"array_prop"},
				Properties: map[string]*base.SchemaProxy{
					"array_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey this is an array, required!",
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
			readResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"array_prop"},
				Properties: map[string]*base.SchemaProxy{
					"array_prop": base.CreateSchemaProxy(&base.Schema{
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
			}),
			want: []ir.ResourceAttribute{
				{
					Name: "array_prop",
					ListNested: &ir.ResourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey this is an array, required!"),
						NestedObject: ir.ResourceAttributeNestedObject{
							Attributes: []ir.ResourceAttribute{
								{
									Name: "float64_prop",
									Float64: &ir.ResourceFloat64Attribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey this is a float64!"),
									},
								},
								{
									Name: "nested_array_prop",
									ListNested: &ir.ResourceListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey this is a nested array, required!"),
										NestedObject: ir.ResourceAttributeNestedObject{
											Attributes: []ir.ResourceAttribute{
												{
													Name: "super_nested_string",
													String: &ir.ResourceStringAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey this is a string!"),
														Sensitive:                pointer(false),
													},
												},
												{
													Name: "super_nested_bool_one",
													Bool: &ir.ResourceBoolAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey this is a boolean!"),
													},
												},
												{
													Name: "super_nested_bool_two",
													Bool: &ir.ResourceBoolAttribute{
														ComputedOptionalRequired: ir.Required,
														Description:              pointer("hey this is a boolean, required!"),
													},
												},
												{
													Name: "super_nested_int64",
													Int64: &ir.ResourceInt64Attribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey this is a integer!"),
													},
												},
											},
										},
									},
								},
								{
									Name: "number_prop",
									Number: &ir.ResourceNumberAttribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey this is a number!"),
									},
								},
							},
						},
					},
				},
			},
		},
		"deep merge list array with object element types": {
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
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
											"deep_nested_float64": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "float",
											}),
											"deep_nested_string": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										},
									}),
								},
							}),
						},
					}),
				},
			}),
			createResponseSchema: base.CreateSchemaProxy(&base.Schema{
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
			}),
			want: []ir.ResourceAttribute{
				{
					Name: "array_prop",
					List: &ir.ResourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey this is an array!"),
						ElementType: ir.ElementType{
							List: &ir.ListElement{
								ElementType: &ir.ElementType{
									Object: []ir.ObjectElement{
										{
											Name: "deep_nested_float64",
											ElementType: &ir.ElementType{
												Float64: &ir.Float64Element{},
											},
										},
										{
											Name: "deep_nested_string",
											ElementType: &ir.ElementType{
												String: &ir.StringElement{},
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
																			Name: "deep_deep_nested_string",
																			ElementType: &ir.ElementType{
																				String: &ir.StringElement{},
																			},
																		},
																		{
																			Name: "deep_deep_nested_bool",
																			ElementType: &ir.ElementType{
																				Bool: &ir.BoolElement{},
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

			mapper := resource.NewResourceMapper(map[string]explorer.Resource{
				"test_resource": {
					CreateOp: createTestCreateOp(testCase.createRequestSchema, testCase.createResponseSchema),
					ReadOp:   createTestReadOp(testCase.readResponseSchema, testCase.readParams),
				},
			}, config.Config{})
			got, err := mapper.MapToIR()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if len(got) != 1 {
				t.Fatalf("expected only one resource, got: %d", len(got))
			}

			if diff := cmp.Diff(got[0].Schema.Attributes, testCase.want); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func createTestCreateOp(request *base.SchemaProxy, response *base.SchemaProxy) *high.Operation {
	return &high.Operation{
		RequestBody: &high.RequestBody{
			Content: map[string]*high.MediaType{
				"application/json": {
					Schema: request,
				},
			},
		},
		Responses: &high.Responses{
			Codes: map[string]*high.Response{
				"201": {
					Content: map[string]*high.MediaType{
						"application/json": {
							Schema: response,
						},
					},
				},
			},
		},
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
