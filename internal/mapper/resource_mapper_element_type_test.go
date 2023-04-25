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
								Type:   []string{"integer"},
								Format: "int64",
							}),
						},
					}),
					"int64_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of int64s, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
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
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of int64s."),
										ElementType: ir.ElementType{
											Int64: &ir.ElementTypeInt64{},
										},
									},
								},
							},
							{
								Name: "int64_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of int64s, required."),
										ElementType: ir.ElementType{
											Int64: &ir.ElementTypeInt64{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"list attributes with number (integer) element type": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"number_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"number_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of numbers.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"integer"},
							}),
						},
					}),
					"number_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of numbers, required.",
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
								Name: "number_list_prop",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of numbers."),
										ElementType: ir.ElementType{
											Number: &ir.ElementTypeNumber{},
										},
									},
								},
							},
							{
								Name: "number_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of numbers, required."),
										ElementType: ir.ElementType{
											Number: &ir.ElementTypeNumber{},
										},
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
				Required: []string{"float64_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"float64_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of float64s.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "double",
							}),
						},
					}),
					"float64_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of float64s, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "double",
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
								Name: "float64_list_prop",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of float64s."),
										ElementType: ir.ElementType{
											Float64: &ir.ElementTypeFloat64{},
										},
									},
								},
							},
							{
								Name: "float64_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of float64s, required."),
										ElementType: ir.ElementType{
											Float64: &ir.ElementTypeFloat64{},
										},
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
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of numbers."),
										ElementType: ir.ElementType{
											Number: &ir.ElementTypeNumber{},
										},
									},
								},
							},
							{
								Name: "number_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of numbers, required."),
										ElementType: ir.ElementType{
											Number: &ir.ElementTypeNumber{},
										},
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
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of strings."),
										ElementType: ir.ElementType{
											String: &ir.ElementTypeString{},
										},
									},
								},
							},
							{
								Name: "string_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of strings, required."),
										ElementType: ir.ElementType{
											String: &ir.ElementTypeString{},
										},
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
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of bools."),
										ElementType: ir.ElementType{
											Bool: &ir.ElementTypeBool{},
										},
									},
								},
							},
							{
								Name: "bool_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of bools, required."),
										ElementType: ir.ElementType{
											Bool: &ir.ElementTypeBool{},
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
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a list of lists."),
										ElementType: ir.ElementType{
											List: &ir.ElementTypeList{
												ElementType: ir.ElementType{
													Object: []ir.ElementTypeObject{
														{
															Name: "float64_prop",
															Type: ir.ElementType{
																Float64: &ir.ElementTypeFloat64{},
															},
														},
														{
															Name: "int64_prop",
															Type: ir.ElementType{
																Int64: &ir.ElementTypeInt64{},
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
								Name: "nested_list_prop_required",
								Type: ir.ResourceAttributeType{
									List: &ir.ResourceList{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list of lists, required."),
										ElementType: ir.ElementType{
											List: &ir.ElementTypeList{
												ElementType: ir.ElementType{
													Object: []ir.ElementTypeObject{
														{
															Name: "bool_prop",
															Type: ir.ElementType{
																Bool: &ir.ElementTypeBool{},
															},
														},
														{
															Name: "string_prop",
															Type: ir.ElementType{
																String: &ir.ElementTypeString{},
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
