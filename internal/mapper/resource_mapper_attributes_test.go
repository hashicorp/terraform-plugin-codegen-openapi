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
						Format:      "int64",
						Description: "hey there! I'm an int64 type.",
					}),
					"int64_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Format:      "int64",
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
								Type: ir.ResourceAttributeType{
									Int64: &ir.ResourceInt64{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm an int64 type."),
									},
								},
							},
							{
								Name: "int64_prop_required",
								Type: ir.ResourceAttributeType{
									Int64: &ir.ResourceInt64{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm an int64 type, required."),
									},
								},
							},
						},
					},
				},
			},
		},
		"number (integer) attributes": {
			requestSchema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"int_number_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"int_number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Description: "hey there! I'm a number type, from an integer.",
					}),
					"int_number_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"integer"},
						Description: "hey there! I'm a number type, from an integer, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "int_number_prop",
								Type: ir.ResourceAttributeType{
									Number: &ir.ResourceNumber{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a number type, from an integer."),
									},
								},
							},
							{
								Name: "int_number_prop_required",
								Type: ir.ResourceAttributeType{
									Number: &ir.ResourceNumber{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a number type, from an integer, required."),
									},
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
				Required: []string{"number_float64_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"number_float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "double",
						Description: "hey there! I'm a float64 type.",
					}),
					"number_float64_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Format:      "double",
						Description: "hey there! I'm a float64 type, required.",
					}),
				},
			},
			expectedIR: &[]ir.Resource{
				{
					Name: "test_resource",
					Schema: ir.ResourceSchema{
						Attributes: []ir.ResourceAttribute{
							{
								Name: "number_float64_prop",
								Type: ir.ResourceAttributeType{
									Float64: &ir.ResourceFloat64{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a float64 type."),
									},
								},
							},
							{
								Name: "number_float64_prop_required",
								Type: ir.ResourceAttributeType{
									Float64: &ir.ResourceFloat64{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a float64 type, required."),
									},
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
								Type: ir.ResourceAttributeType{
									Number: &ir.ResourceNumber{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a number type."),
									},
								},
							},
							{
								Name: "number_prop_required",
								Type: ir.ResourceAttributeType{
									Number: &ir.ResourceNumber{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a number type, required."),
									},
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
								Type: ir.ResourceAttributeType{
									String: &ir.ResourceString{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a string type, not sensitive, required."),
										Sensitive:                pointer(false),
									},
								},
							},
							{
								Name: "string_sensitive_prop",
								Type: ir.ResourceAttributeType{
									String: &ir.ResourceString{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a string type, sensitive"),
										Sensitive:                pointer(true),
									},
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
								Type: ir.ResourceAttributeType{
									Bool: &ir.ResourceBool{
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a bool type."),
									},
								},
							},
							{
								Name: "bool_prop_required",
								Type: ir.ResourceAttributeType{
									Bool: &ir.ResourceBool{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a bool type, required."),
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
								Type: ir.ResourceAttributeType{
									SingleNested: &ir.SingleNestedAttribute{
										Attributes: []ir.ResourceAttribute{
											{
												Name: "nested_obj_prop_required",
												Type: ir.ResourceAttributeType{
													SingleNested: &ir.SingleNestedAttribute{
														Attributes: []ir.ResourceAttribute{
															{
																Name: "nested_float64",
																Type: ir.ResourceAttributeType{
																	Float64: &ir.ResourceFloat64{
																		ComputedOptionalRequired: ir.Optional,
																		Description:              pointer("hey there! I'm a nested float64 type."),
																	},
																},
															},
															{
																Name: "nested_int64_required",
																Type: ir.ResourceAttributeType{
																	Int64: &ir.ResourceInt64{
																		ComputedOptionalRequired: ir.Required,
																		Description:              pointer("hey there! I'm a nested int64 type, required."),
																	},
																},
															},
														},
														ComputedOptionalRequired: ir.Required,
														Description:              pointer("hey there! I'm a single nested object type, required."),
													},
												},
											},
										},
										ComputedOptionalRequired: ir.Optional,
										Description:              pointer("hey there! I'm a single nested object type."),
									},
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
								Type: ir.ResourceAttributeType{
									ListNested: &ir.ListNestedAttribute{
										ComputedOptionalRequired: ir.Required,
										Description:              pointer("hey there! I'm a list nested array type, required."),
										NestedObject: ir.NestedObjectClass{
											Attributes: []ir.ResourceAttribute{
												{
													Name: "nested_float64",
													Type: ir.ResourceAttributeType{
														Float64: &ir.ResourceFloat64{
															ComputedOptionalRequired: ir.Optional,
															Description:              pointer("hey there! I'm a nested float64 type."),
														},
													},
												},
												{
													Name: "nested_int64_required",
													Type: ir.ResourceAttributeType{
														Int64: &ir.ResourceInt64{
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
