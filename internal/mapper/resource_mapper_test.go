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
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func TestResourceMapper_basic_merges(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		createRequestSchema  *base.SchemaProxy
		createResponseSchema *base.SchemaProxy
		readResponseSchema   *base.SchemaProxy
		readParams           []*high.Parameter
		schemaOptions        explorer.SchemaOptions
		want                 resource.Attributes
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
					Name:        "string_prop",
					In:          "path",
					Description: "hey this is a string, overridden!",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      util.OAS_format_password,
						Description: "you shouldn't see this because the description is overridden!",
					}),
				},
			},
			want: resource.Attributes{
				{
					Name: "bool_prop",
					Bool: &resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is a bool, required!"),
					},
				},
				{
					Name: "int64_prop",
					Int64: &resource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is an int64, required!"),
					},
				},
				{
					Name: "number_prop",
					Number: &resource.NumberAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("hey this is a number!"),
					},
				},
				{
					Name: "string_prop",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is a string, overridden!"),
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
										Description: "hey this is a integer, switched to computed optional!",
									}),
								},
							}),
						},
					}),
				},
			},
			want: resource.Attributes{
				{
					Name: "nested_object_one",
					SingleNested: &resource.SingleNestedAttribute{
						Attributes: []resource.Attribute{
							{
								Name: "bool_prop",
								Bool: &resource.BoolAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("hey this is a bool, required!"),
								},
							},
							{
								Name: "nested_object_two",
								SingleNested: &resource.SingleNestedAttribute{
									Attributes: []resource.Attribute{
										{
											Name: "bool_prop",
											Bool: &resource.BoolAttribute{
												ComputedOptionalRequired: schema.Computed,
												Description:              pointer("hey this is a bool!"),
											},
										},
										{
											Name: "number_prop",
											Number: &resource.NumberAttribute{
												ComputedOptionalRequired: schema.Computed,
												Description:              pointer("hey this is a number!"),
											},
										},
										{
											Name: "int64_prop",
											Int64: &resource.Int64Attribute{
												ComputedOptionalRequired: schema.ComputedOptional,
												Description:              pointer("hey this is a integer, switched to computed optional!"),
											},
										},
									},
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("hey this is an object!"),
								},
							},
							{
								Name: "string_prop",
								String: &resource.StringAttribute{
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("hey this is a string!"),
									Sensitive:                pointer(true),
								},
							},
						},
						ComputedOptionalRequired: schema.Required,
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
														Description: "hey this is a boolean, switched to computed optional!",
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
			want: resource.Attributes{
				{
					Name: "array_prop",
					ListNested: &resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey this is an array, required!"),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "float64_prop",
									Float64: &resource.Float64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey this is a float64!"),
									},
								},
								{
									Name: "nested_array_prop",
									ListNested: &resource.ListNestedAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey this is a nested array, required!"),
										NestedObject: resource.NestedAttributeObject{
											Attributes: []resource.Attribute{
												{
													Name: "super_nested_string",
													String: &resource.StringAttribute{
														ComputedOptionalRequired: schema.ComputedOptional,
														Description:              pointer("hey this is a string!"),
													},
												},
												{
													Name: "super_nested_bool_one",
													Bool: &resource.BoolAttribute{
														ComputedOptionalRequired: schema.Computed,
														Description:              pointer("hey this is a boolean!"),
													},
												},
												{
													Name: "super_nested_bool_two",
													Bool: &resource.BoolAttribute{
														ComputedOptionalRequired: schema.Computed,
														Description:              pointer("hey this is a boolean, switched to computed optional!"),
													},
												},
												{
													Name: "super_nested_int64",
													Int64: &resource.Int64Attribute{
														ComputedOptionalRequired: schema.Computed,
														Description:              pointer("hey this is a integer!"),
													},
												},
											},
										},
									},
								},
								{
									Name: "number_prop",
									Number: &resource.NumberAttribute{
										ComputedOptionalRequired: schema.Computed,
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
			want: resource.Attributes{
				{
					Name: "array_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is an array!"),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "deep_nested_float64",
												Float64: &schema.Float64Type{},
											},
											{
												Name:   "deep_nested_string",
												String: &schema.StringType{},
											},
											{
												Name: "deep_nested_bool",
												Bool: &schema.BoolType{},
											},
											{
												Name:  "deep_nested_int64",
												Int64: &schema.Int64Type{},
											},
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
																				Name:   "deep_deep_nested_string",
																				String: &schema.StringType{},
																			},
																			{
																				Name: "deep_deep_nested_bool",
																				Bool: &schema.BoolType{},
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
		"precedence and configurability": {
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Required: []string{
					"create_request_required_create_request_only",
					"create_request_required_create_response",
					"create_request_required_read_parameter_optional",
					"create_request_required_read_parameter_required",
					"create_request_required_read_response",
				},
				Properties: map[string]*base.SchemaProxy{
					"create_request_optional_create_request_only": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_optional_create_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_optional_read_parameter_optional": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_optional_read_parameter_required": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_optional_read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_required_create_request_only": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_required_create_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_required_read_parameter_optional": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_required_read_parameter_required": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_request_required_read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
			createResponseSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					// Simulate API returning parameter in response
					"create_request_optional_create_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					// Simulate API returning parameter in response
					"create_request_required_create_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"create_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
			readParams: []*high.Parameter{
				{
					Name:     "create_request_optional_read_parameter_optional",
					Required: false,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "create_request_optional_read_parameter_required",
					Required: true,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "create_request_required_read_parameter_optional",
					Required: false,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "create_request_required_read_parameter_required",
					Required: true,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				// Edge case where new read request properties do not align with
				// other request/response properties. Provider developers would
				// want to configure the converter to remap existing properties
				// to align with these or risk the converter returning API-level
				// details to practitioners, which are conventionally hidden.
				{
					Name:     "read_parameter_optional",
					Required: false,
					In:       "path",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
				{
					Name:     "read_parameter_required",
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
					"create_request_optional_read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					// Simulate API returning parameter in response
					"create_request_required_read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
					"read_response": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
			want: resource.Attributes{
				{
					Name: "create_request_optional_create_request_only",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "create_request_optional_create_response",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "create_request_optional_read_parameter_optional",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "create_request_optional_read_parameter_required",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "create_request_optional_read_response",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "create_request_required_create_request_only",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "create_request_required_create_response",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "create_request_required_read_parameter_optional",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "create_request_required_read_parameter_required",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "create_request_required_read_response",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "create_response",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Computed,
					},
				},
				{
					Name: "read_response",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Computed,
					},
				},
				{
					Name: "read_parameter_optional",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
				{
					Name: "read_parameter_required",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
					},
				},
			},
		},
		"parameter match for path and query params": {
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
				Type:     []string{"object"},
				Required: []string{"attribute_required"},
				Properties: map[string]*base.SchemaProxy{
					"attribute_required": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
					}),
				},
			}),
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
					"attribute_computed": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"boolean"},
					}),
				},
			}),
			schemaOptions: explorer.SchemaOptions{
				AttributeOptions: explorer.AttributeOptions{
					Aliases: map[string]string{
						"read_path_parameter":  "attribute_required",
						"read_query_parameter": "attribute_computed",
					},
				},
			},
			want: resource.Attributes{
				{
					Name: "attribute_required",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
				{
					Name: "attribute_computed",
					Bool: &resource.BoolAttribute{
						ComputedOptionalRequired: schema.Computed,
					},
				},
			},
		},
		"ignore bool prop across all ops": {
			schemaOptions: explorer.SchemaOptions{
				Ignores: []string{
					"bool_prop",
					"nested_obj.bool_prop",
					"nested_array.deep_nested_bool",
					"nested_map.deep_nested_bool",
				},
			},
			createRequestSchema: base.CreateSchemaProxy(&base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "This boolean is going to be ignored!",
					}),
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey this is a number!",
					}),
					"nested_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey this is a map!",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								Properties: map[string]*base.SchemaProxy{
									"deep_nested_bool": base.CreateSchemaProxy(&base.Schema{
										Type: []string{"boolean"},
									}),
									"deep_nested_int64": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"integer"},
										Description: "hey this is an int64!",
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
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "This boolean is going to be ignored!",
					}),
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey this is a number!",
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
										Properties: map[string]*base.SchemaProxy{
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
						Properties: map[string]*base.SchemaProxy{
							"bool_prop": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"boolean"},
								Description: "This boolean is going to be ignored!",
							}),
							"string_prop": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey this is a string!",
							}),
						},
					}),
				},
			}),
			readParams: []*high.Parameter{
				{
					Name:     "bool_prop",
					Required: true,
					In:       "query",
					Schema: base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "This boolean is going to be ignored!",
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
			want: resource.Attributes{
				{
					Name: "nested_map",
					MapNested: &resource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is a map!"),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "deep_nested_int64",
									Int64: &resource.Int64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey this is an int64!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "number_prop",
					Number: &resource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is a number!"),
					},
				},
				{
					Name: "nested_array",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("hey this is an array!"),
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
					Name: "nested_obj",
					SingleNested: &resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Computed,
						Attributes: []resource.Attribute{
							{
								Name: "string_prop",
								String: &resource.StringAttribute{
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("hey this is a string!"),
								},
							},
						},
					},
				},
				{
					Name: "float64_prop",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey this is a float64!"),
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mapper := mapper.NewResourceMapper(map[string]explorer.Resource{
				"test_resource": {
					CreateOp:      createTestCreateOp(testCase.createRequestSchema, testCase.createResponseSchema),
					ReadOp:        createTestReadOp(testCase.readResponseSchema, testCase.readParams),
					SchemaOptions: testCase.schemaOptions,
				},
			}, config.Config{})
			got, err := mapper.MapToIR(slog.Default())
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
