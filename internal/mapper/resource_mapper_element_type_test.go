package mapper_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestResourceMapper_ListAttributes_PrimitiveElements(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestSchema *base.Schema
		expectedIR    *[]ir.Resource
	}{
		"list attributes with int64 element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"int64_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"int64_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of int64s.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"integer"},
							}),
						},
					}),
					"int64_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of int64s, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"integer"},
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
								Name: "int64_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of int64s."),
									ElementType: ir.ElementType{
										Int64: &ir.Int64Element{},
									},
								},
							},
							{
								Name: "int64_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of int64s, required."),
									ElementType: ir.ElementType{
										Int64: &ir.Int64Element{},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with float64 element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_list_prop_required", "float_float64_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"double_float64_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of float64s.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "double",
							}),
						},
					}),
					"double_float64_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of float64s, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "double",
							}),
						},
					}),
					"float_float64_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of float64s.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "float",
							}),
						},
					}),
					"float_float64_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of float64s, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "float",
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
								Name: "double_float64_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of float64s."),
									ElementType: ir.ElementType{
										Float64: &ir.Float64Element{},
									},
								},
							},
							{
								Name: "double_float64_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of float64s, required."),
									ElementType: ir.ElementType{
										Float64: &ir.Float64Element{},
									},
								},
							},
							{
								Name: "float_float64_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of float64s."),
									ElementType: ir.ElementType{
										Float64: &ir.Float64Element{},
									},
								},
							},
							{
								Name: "float_float64_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of float64s, required."),
									ElementType: ir.ElementType{
										Float64: &ir.Float64Element{},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with number element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"number_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"number_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of numbers.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"number"},
							}),
						},
					}),
					"number_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of numbers, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"number"},
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
								Name: "number_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of numbers."),
									ElementType: ir.ElementType{
										Number: &ir.NumberElement{},
									},
								},
							},
							{
								Name: "number_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of numbers, required."),
									ElementType: ir.ElementType{
										Number: &ir.NumberElement{},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with string element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string"},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string"},
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
								Name: "string_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of strings."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
							{
								Name: "string_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of strings, required."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with bool element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
					"bool_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
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
								Name: "bool_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of bools."),
									ElementType: ir.ElementType{
										Bool: &ir.BoolElement{},
									},
								},
							},
							{
								Name: "bool_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of bools, required."),
									ElementType: ir.ElementType{
										Bool: &ir.BoolElement{},
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

func TestResourceMapper_ListAttributes_NestedElements(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestSchema *base.Schema
		expectedIR    *[]ir.Resource
	}{
		"list attributes with list and nested object element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of lists.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										},
									}),
								},
							}),
						},
					}),
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of lists, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										},
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
								Name: "nested_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of lists."),
									ElementType: ir.ElementType{
										List: &ir.ListElement{
											ElementType: &ir.ElementType{
												Object: []ir.ObjectElement{
													{
														Name: "float64_prop",
														ElementType: &ir.ElementType{
															Float64: &ir.Float64Element{},
														},
													},
													{
														Name: "int64_prop",
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
							{
								Name: "nested_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of lists, required."),
									ElementType: ir.ElementType{
										List: &ir.ListElement{
											ElementType: &ir.ElementType{
												Object: []ir.ObjectElement{
													{
														Name: "bool_prop",
														ElementType: &ir.ElementType{
															Bool: &ir.BoolElement{},
														},
													},
													{
														Name: "string_prop",
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

func TestResourceMapper_ListAttributes_NullableMultiTypes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		requestSchema *base.Schema
		expectedIR    *[]ir.Resource
	}{
		"list attributes with nullable element type - Type array": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"null", "string"},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string", "null"},
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
								Name: "string_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of nullable strings."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
							{
								Name: "string_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of nullable strings, required."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with nullable element type - anyOf": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								AnyOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
								},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								AnyOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
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
								Name: "string_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of nullable strings."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
							{
								Name: "string_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of nullable strings, required."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with nullable element type - oneOf": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"string_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								OneOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
								},
							}),
						},
					}),
					"string_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of nullable strings, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								OneOf: []*base.SchemaProxy{
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"string"},
									}),
									base.CreateSchemaProxy(&base.Schema{
										Type: []string{"null"},
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
								Name: "string_list_prop",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.ComputedOptional,
									Description:              pointer("hey there! I'm a list of nullable strings."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
									},
								},
							},
							{
								Name: "string_list_prop_required",
								List: &ir.ResourceListAttribute{
									ComputedOptionalRequired: ir.Required,
									Description:              pointer("hey there! I'm a list of nullable strings, required."),
									ElementType: ir.ElementType{
										String: &ir.StringElement{},
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
