// Should eventually move these tests to a black-box test
package mapper

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeepMerge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		main   []ir.ResourceAttribute
		second []ir.ResourceAttribute
		third  []ir.ResourceAttribute
		want   []ir.ResourceAttribute
	}{
		"merge three basic slices": {
			main: []ir.ResourceAttribute{
				{
					Name: "test_attribute_one",
					Bool: &ir.ResourceBoolAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_attribute_one!"),
					},
				},
			},
			second: []ir.ResourceAttribute{
				// This one will be skipped because it exists in the main slice
				{
					Name: "test_attribute_one",
					Bool: &ir.ResourceBoolAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("this one should be skipped!"),
					},
				},
				{
					Name: "test_attribute_two",
					Int64: &ir.ResourceInt64Attribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there test_attribute_two!"),
					},
				},
			},
			third: []ir.ResourceAttribute{
				{
					Name: "test_attribute_third",
					String: &ir.ResourceStringAttribute{
						ComputedOptionalRequired: ir.Computed,
						Description:              pointer("hey there test_attribute_third!"),
					},
				},
			},
			want: []ir.ResourceAttribute{
				{
					Name: "test_attribute_one",
					Bool: &ir.ResourceBoolAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_attribute_one!"),
					},
				},
				{
					Name: "test_attribute_two",
					Int64: &ir.ResourceInt64Attribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there test_attribute_two!"),
					},
				},
				{
					Name: "test_attribute_third",
					String: &ir.ResourceStringAttribute{
						ComputedOptionalRequired: ir.Computed,
						Description:              pointer("hey there test_attribute_third!"),
					},
				},
			},
		},
		"merge three deep nested slices - single object": {
			main: []ir.ResourceAttribute{
				{
					Name: "test_nested_attribute",
					SingleNested: &ir.ResourceSingleNestedAttribute{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "test_attribute_nested_one",
								Bool: &ir.ResourceBoolAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there test_attribute_nested_one!"),
								},
							},
						},
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_nested_attribute!"),
					},
				},
			},
			second: []ir.ResourceAttribute{
				{
					Name: "test_nested_attribute",
					SingleNested: &ir.ResourceSingleNestedAttribute{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "test_attribute_nested_two",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there test_attribute_nested_two!"),
								},
							},
							{
								Name: "test_attribute_deep_nested",
								SingleNested: &ir.ResourceSingleNestedAttribute{
									Attributes: []ir.ResourceAttribute{
										{
											Name: "test_attribute_deep_nested_one",
											Bool: &ir.ResourceBoolAttribute{
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey there test_attribute_deep_nested_one!"),
											},
										},
										{
											Name: "test_attribute_deep_nested_three",
											Int64: &ir.ResourceInt64Attribute{
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey there test_attribute_deep_nested_three!"),
											},
										},
									},
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there test_attribute_deep_nested!"),
								},
							},
						},
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_nested_attribute!"),
					},
				},
			},
			third: []ir.ResourceAttribute{
				{
					Name: "test_nested_attribute",
					SingleNested: &ir.ResourceSingleNestedAttribute{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "test_attribute_nested_three",
								Float64: &ir.ResourceFloat64Attribute{
									ComputedOptionalRequired: ir.Computed,
									Description:              pointer("hey there test_attribute_nested_three!"),
								},
							},
							{
								Name: "test_attribute_deep_nested",
								SingleNested: &ir.ResourceSingleNestedAttribute{
									Attributes: []ir.ResourceAttribute{
										{
											Name: "test_attribute_deep_nested_two",
											String: &ir.ResourceStringAttribute{
												ComputedOptionalRequired: ir.ComputedOptional,
												Description:              pointer("hey there test_attribute_deep_nested_two!"),
											},
										},
										{
											Name: "test_attribute_deep_nested_three",
											Float64: &ir.ResourceFloat64Attribute{
												ComputedOptionalRequired: ir.Computed,
												Description:              pointer("hey there test_attribute_deep_nested_three!"),
											},
										},
									},
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there test_attribute_deep_nested!"),
								},
							},
						},
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_nested_attribute!"),
					},
				},
			},
			want: []ir.ResourceAttribute{
				{
					Name: "test_nested_attribute",
					SingleNested: &ir.ResourceSingleNestedAttribute{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "test_attribute_nested_one",
								Bool: &ir.ResourceBoolAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there test_attribute_nested_one!"),
								},
							},
							{
								Name: "test_attribute_nested_two",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there test_attribute_nested_two!"),
								},
							},
							{
								Name: "test_attribute_deep_nested",
								SingleNested: &ir.ResourceSingleNestedAttribute{
									Attributes: []ir.ResourceAttribute{
										{
											Name: "test_attribute_deep_nested_one",
											Bool: &ir.ResourceBoolAttribute{
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey there test_attribute_deep_nested_one!"),
											},
										},
										{
											Name: "test_attribute_deep_nested_three",
											Int64: &ir.ResourceInt64Attribute{
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey there test_attribute_deep_nested_three!"),
											},
										},
										{
											Name: "test_attribute_deep_nested_two",
											String: &ir.ResourceStringAttribute{
												ComputedOptionalRequired: ir.ComputedOptional,
												Description:              pointer("hey there test_attribute_deep_nested_two!"),
											},
										},
									},
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there test_attribute_deep_nested!"),
								},
							},
							{
								Name: "test_attribute_nested_three",
								Float64: &ir.ResourceFloat64Attribute{
									ComputedOptionalRequired: ir.Computed,
									Description:              pointer("hey there test_attribute_nested_three!"),
								},
							},
						},
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_nested_attribute!"),
					},
				},
			},
		},
		"merge three deep nested slices - list nested array": {
			main: []ir.ResourceAttribute{
				{
					Name: "test_listnested_attribute",
					ListNested: &ir.ResourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_listnested_attribute!"),
						NestedObject: ir.ResourceAttributeNestedObject{
							Attributes: []ir.ResourceAttribute{
								{
									Name: "nested_listnested_attribute",
									ListNested: &ir.ResourceListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there nested_listnested_attribute!"),
										NestedObject: ir.ResourceAttributeNestedObject{
											Attributes: []ir.ResourceAttribute{
												{
													Name: "super_nested_string_attribute",
													String: &ir.ResourceStringAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey there super_nested_string_attribute!"),
													},
												},
											},
										},
									},
								},
								{
									Name: "nested_float_attribute",
									Float64: &ir.ResourceFloat64Attribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey there nested_float_attribute!"),
									},
								},
							},
						},
					},
				},
			},
			second: []ir.ResourceAttribute{
				{
					Name: "test_listnested_attribute",
					ListNested: &ir.ResourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_listnested_attribute!"),
						NestedObject: ir.ResourceAttributeNestedObject{
							Attributes: []ir.ResourceAttribute{
								{
									Name: "nested_listnested_attribute",
									ListNested: &ir.ResourceListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there nested_listnested_attribute!"),
										NestedObject: ir.ResourceAttributeNestedObject{
											Attributes: []ir.ResourceAttribute{
												{
													Name: "super_nested_bool_attribute_one",
													Bool: &ir.ResourceBoolAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey there super_nested_bool_attribute_one!"),
													},
												},
												{
													Name: "super_nested_bool_attribute_two",
													Bool: &ir.ResourceBoolAttribute{
														ComputedOptionalRequired: ir.Required,
														Description:              pointer("hey there super_nested_bool_attribute_two!"),
													},
												},
											},
										},
									},
								},
								{
									Name: "nested_float_attribute",
									Float64: &ir.ResourceFloat64Attribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey there nested_float_attribute!"),
									},
								},
								{
									Name: "nested_number_attribute",
									Number: &ir.ResourceNumberAttribute{
										ComputedOptionalRequired: ir.Computed,
										Description:              pointer("hey there nested_number_attribute!"),
									},
								},
							},
						},
					},
				},
			},
			third: []ir.ResourceAttribute{
				{
					Name: "test_listnested_attribute",
					ListNested: &ir.ResourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_listnested_attribute!"),
						NestedObject: ir.ResourceAttributeNestedObject{
							Attributes: []ir.ResourceAttribute{
								{
									Name: "nested_listnested_attribute",
									ListNested: &ir.ResourceListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there nested_listnested_attribute!"),
										NestedObject: ir.ResourceAttributeNestedObject{
											Attributes: []ir.ResourceAttribute{
												{
													Name: "super_nested_int64_attribute",
													Int64: &ir.ResourceInt64Attribute{
														ComputedOptionalRequired: ir.Computed,
														Description:              pointer("hey there super_nested_int64_attribute!"),
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
			want: []ir.ResourceAttribute{
				{
					Name: "test_listnested_attribute",
					ListNested: &ir.ResourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_listnested_attribute!"),
						NestedObject: ir.ResourceAttributeNestedObject{
							Attributes: []ir.ResourceAttribute{
								{
									Name: "nested_listnested_attribute",
									ListNested: &ir.ResourceListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there nested_listnested_attribute!"),
										NestedObject: ir.ResourceAttributeNestedObject{
											Attributes: []ir.ResourceAttribute{
												{
													Name: "super_nested_string_attribute",
													String: &ir.ResourceStringAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey there super_nested_string_attribute!"),
													},
												},
												{
													Name: "super_nested_bool_attribute_one",
													Bool: &ir.ResourceBoolAttribute{
														ComputedOptionalRequired: ir.ComputedOptional,
														Description:              pointer("hey there super_nested_bool_attribute_one!"),
													},
												},
												{
													Name: "super_nested_bool_attribute_two",
													Bool: &ir.ResourceBoolAttribute{
														ComputedOptionalRequired: ir.Required,
														Description:              pointer("hey there super_nested_bool_attribute_two!"),
													},
												},
												{
													Name: "super_nested_int64_attribute",
													Int64: &ir.ResourceInt64Attribute{
														ComputedOptionalRequired: ir.Computed,
														Description:              pointer("hey there super_nested_int64_attribute!"),
													},
												},
											},
										},
									},
								},
								{
									Name: "nested_float_attribute",
									Float64: &ir.ResourceFloat64Attribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey there nested_float_attribute!"),
									},
								},
								{
									Name: "nested_number_attribute",
									Number: &ir.ResourceNumberAttribute{
										ComputedOptionalRequired: ir.Computed,
										Description:              pointer("hey there nested_number_attribute!"),
									},
								},
							},
						},
					},
				},
			},
		},
		"merge three deep nested slices - lists with object element types": {
			main: []ir.ResourceAttribute{
				{
					Name: "test_list_attribute",
					List: &ir.ResourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_list_attribute!"),
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
			second: []ir.ResourceAttribute{
				{
					Name: "test_list_attribute",
					List: &ir.ResourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_list_attribute!"),
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
											Name: "deep_nested_bool",
											ElementType: &ir.ElementType{
												Bool: &ir.BoolElement{},
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
			third: []ir.ResourceAttribute{
				{
					Name: "test_list_attribute",
					List: &ir.ResourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_list_attribute!"),
						ElementType: ir.ElementType{
							List: &ir.ListElement{
								ElementType: &ir.ElementType{
									Object: []ir.ObjectElement{
										{
											Name: "deep_nested_int64",
											ElementType: &ir.ElementType{
												Float64: &ir.Float64Element{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []ir.ResourceAttribute{
				{
					Name: "test_list_attribute",
					List: &ir.ResourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there test_list_attribute!"),
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
										{
											Name: "deep_nested_bool",
											ElementType: &ir.ElementType{
												Bool: &ir.BoolElement{},
											},
										},
										{
											Name: "deep_nested_int64",
											ElementType: &ir.ElementType{
												Float64: &ir.Float64Element{},
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

			got := deepMergeAttributes(testCase.main, testCase.second, testCase.third)

			if diff := cmp.Diff(got, &testCase.want); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func pointer(str string) *string {
	return &str
}
