package schema_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: add error tests

func TestBuildListResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.ResourceAttribute
	}{
		"list nested attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"list attributes with list and nested object element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := schema.OASSchema{Schema: testCase.schema}
			attributes, err := schema.BuildResourceAttributes()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(attributes, testCase.expectedAttributes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestBuildListDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.DataSourceAttribute
	}{
		"list nested attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "nested_list_prop_required",
					ListNested: &ir.DataSourceListNestedAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list nested array type, required."),
						NestedObject: ir.DataSourceAttributeNestedObject{
							Attributes: []ir.DataSourceAttribute{
								{
									Name: "nested_float64",
									Float64: &ir.DataSourceFloat64Attribute{
										ComputedOptionalRequired: ir.ComputedOptional,
										Description:              pointer("hey there! I'm a nested float64 type."),
									},
								},
								{
									Name: "nested_int64_required",
									Int64: &ir.DataSourceInt64Attribute{
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
		"list attributes with list and nested object element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "nested_list_prop",
					List: &ir.DataSourceListAttribute{
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
					List: &ir.DataSourceListAttribute{
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
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := schema.OASSchema{Schema: testCase.schema}
			attributes, err := schema.BuildDataSourceAttributes()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(attributes, testCase.expectedAttributes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
