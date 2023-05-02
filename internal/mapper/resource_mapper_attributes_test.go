package mapper_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestResourceMapper_PrimitiveAttributes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestSchema *base.Schema
		expectedIR    *[]ir.Resource
	}{
		"int64 attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"int64_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"int64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Description: "hey there! I'm an int64 type.",
					}),
					"int64_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Description: "hey there! I'm an int64 type, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "int64_prop",
								Int64: &ir.ResourceInt64Attribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm an int64 type."),
								},
							},
							{
								Name: "int64_prop_required",
								Int64: &ir.ResourceInt64Attribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm an int64 type, required."),
								},
							},
						},
					},
				},
			},
		},
		"float64 attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_prop_required", "float_float64_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"double_float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "double",
						Description: "hey there! I'm a float64 type, from a double.",
					}),
					"double_float64_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "double",
						Description: "hey there! I'm a float64 type, from a double, required.",
					}),
					"float_float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey there! I'm a float64 type, from a float.",
					}),
					"float_float64_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "float",
						Description: "hey there! I'm a float64 type, from a float, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "double_float64_prop",
								Float64: &ir.ResourceFloat64Attribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a float64 type, from a double."),
								},
							},
							{
								Name: "double_float64_prop_required",
								Float64: &ir.ResourceFloat64Attribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a float64 type, from a double, required."),
								},
							},
							{
								Name: "float_float64_prop",
								Float64: &ir.ResourceFloat64Attribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a float64 type, from a float."),
								},
							},
							{
								Name: "float_float64_prop_required",
								Float64: &ir.ResourceFloat64Attribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a float64 type, from a float, required."),
								},
							},
						},
					},
				},
			},
		},
		"number attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"number_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type.",
					}),
					"number_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "number_prop",
								Number: &ir.ResourceNumberAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a number type."),
								},
							},
							{
								Name: "number_prop_required",
								Number: &ir.ResourceNumberAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a number type, required."),
								},
							},
						},
					},
				},
			},
		},
		"string attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Description: "hey there! I'm a string type, not sensitive, required.",
					}),
					"string_sensitive_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      "password",
						Description: "hey there! I'm a string type, sensitive",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "string_prop",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a string type, not sensitive, required."),
									Sensitive:                pointer(false),
								},
							},
							{
								Name: "string_sensitive_prop",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a string type, sensitive"),
									Sensitive:                pointer(true),
								},
							},
						},
					},
				},
			},
		},
		"boolean attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type.",
					}),
					"bool_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "bool_prop",
								Bool: &ir.ResourceBoolAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a bool type."),
								},
							},
							{
								Name: "bool_prop_required",
								Bool: &ir.ResourceBoolAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a bool type, required."),
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
			testResources := createTestResources(testCase.requestSchema)

			resourceMapper := mapper.NewResourceMapper()
			irResources, err := resourceMapper.MapToIR(testResources)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(irResources, testCase.expectedIR); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceMapper_NestedAttributes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestSchema *base.Schema
		expectedIR    *[]ir.Resource
	}{
		"single nested attributes": {
			requestSchema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
								Properties: map[string]*base.SchemaProxy{
									"nested_float64": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"number"},
										Format:      "double",
										Description: "hey there! I'm a nested float64 type.",
									}),
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"integer"},
										Format:      "int64",
										Description: "hey there! I'm a nested int64 type, required.",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "nested_obj_prop",
								SingleNested: &ir.ResourceSingleNestedAttribute{
									Attributes: []ir.ResourceAttribute{
										{
											Name: "nested_obj_prop_required",
											SingleNested: &ir.ResourceSingleNestedAttribute{
												Attributes: []ir.ResourceAttribute{
													{
														Name: "nested_float64",
														Float64: &ir.ResourceFloat64Attribute{
															ComputedOptionalRequired: ir.ComputedOptional,
															Description:              pointer("hey there! I'm a nested float64 type."),
														},
													},
													{
														Name: "nested_int64_required",
														Int64: &ir.ResourceInt64Attribute{
															ComputedOptionalRequired: ir.Required,
															Description:              pointer("hey there! I'm a nested int64 type, required."),
														},
													},
												},
												ComputedOptionalRequired: ir.Required,
												Description:              pointer("hey there! I'm a single nested object type, required."),
											},
										},
									},
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a single nested object type."),
								},
							},
						},
					},
				},
			},
		},
		"list nested attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_float64": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"number"},
										Format:      "double",
										Description: "hey there! I'm a nested float64 type.",
									}),
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"integer"},
										Format:      "int64",
										Description: "hey there! I'm a nested int64 type, required.",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "nested_list_prop_required",
								ListNested: &ir.ResourceListNestedAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list nested array type, required."),
									NestedObject: ir.ResourceAttributeNestedObject{
										Attributes: []ir.ResourceAttribute{
											{
												Name: "nested_float64",
												Float64: &ir.ResourceFloat64Attribute{
													ComputedOptionalRequired: ir.ComputedOptional,
													Description:              pointer("hey there! I'm a nested float64 type."),
												},
											},
											{
												Name: "nested_int64_required",
												Int64: &ir.ResourceInt64Attribute{
													ComputedOptionalRequired: ir.Required,
													Description:              pointer("hey there! I'm a nested int64 type, required."),
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
			testResources := createTestResources(testCase.requestSchema)

			resourceMapper := mapper.NewResourceMapper()
			irResources, err := resourceMapper.MapToIR(testResources)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(irResources, testCase.expectedIR); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResourceMapper_NullableMultiTypes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestSchema *base.Schema
		expectedIR    *[]ir.Resource
	}{
		"nullable type - Type array": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"null", "string"},
						Description: "hey there! I'm a nullable string type.",
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string", "null"},
						Description: "hey there! I'm a nullable string type, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "nullable_string_one",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a nullable string type."),
									Sensitive:                pointer(false),
								},
							},
							{
								Name: "nullable_string_two",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a nullable string type, required."),
									Sensitive:                pointer(false),
								},
							},
						},
					},
				},
			},
		},
		"nullable type - anyOf": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type.",
							}),
						},
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						AnyOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type, required.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
						},
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "nullable_string_one",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a string type."),
									Sensitive:                pointer(false),
								},
							},
							{
								Name: "nullable_string_two",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a string type, required."),
									Sensitive:                pointer(false),
								},
							},
						},
					},
				},
			},
		},
		"nullable type - oneOf": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nullable_string_two"},
				Properties: map[string]*base.SchemaProxy{
					"nullable_string_one": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type.",
							}),
						},
					}),
					"nullable_string_two": base.CreateSchemaProxy(&base.Schema{
						OneOf: []*base.SchemaProxy{
							base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"string"},
								Description: "hey there! I'm a string type, required.",
							}),
							base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null"},
							}),
						},
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "nullable_string_one",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a string type."),
									Sensitive:                pointer(false),
								},
							},
							{
								Name: "nullable_string_two",
								String: &ir.ResourceStringAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a string type, required."),
									Sensitive:                pointer(false),
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
			testResources := createTestResources(testCase.requestSchema)

			resourceMapper := mapper.NewResourceMapper()
			irResources, err := resourceMapper.MapToIR(testResources)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(irResources, testCase.expectedIR); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
