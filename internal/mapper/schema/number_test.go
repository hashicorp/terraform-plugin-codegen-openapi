package schema_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestBuildNumberResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.ResourceAttribute
	}{
		"float64 attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"number attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"list attributes with float64 element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"list attributes with number element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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

func TestBuildNumberDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.DataSourceAttribute
	}{
		"float64 attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "double_float64_prop",
					Float64: &ir.DataSourceFloat64Attribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				{
					Name: "double_float64_prop_required",
					Float64: &ir.DataSourceFloat64Attribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				{
					Name: "float_float64_prop",
					Float64: &ir.DataSourceFloat64Attribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				{
					Name: "float_float64_prop_required",
					Float64: &ir.DataSourceFloat64Attribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"number attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "number_prop",
					Number: &ir.DataSourceNumberAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a number type."),
					},
				},
				{
					Name: "number_prop_required",
					Number: &ir.DataSourceNumberAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a number type, required."),
					},
				},
			},
		},
		"list attributes with float64 element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "double_float64_list_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: ir.ElementType{
							Float64: &ir.Float64Element{},
						},
					},
				},
				{
					Name: "double_float64_list_prop_required",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: ir.ElementType{
							Float64: &ir.Float64Element{},
						},
					},
				},
				{
					Name: "float_float64_list_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: ir.ElementType{
							Float64: &ir.Float64Element{},
						},
					},
				},
				{
					Name: "float_float64_list_prop_required",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: ir.ElementType{
							Float64: &ir.Float64Element{},
						},
					},
				},
			},
		},
		"list attributes with number element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "number_list_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a list of numbers."),
						ElementType: ir.ElementType{
							Number: &ir.NumberElement{},
						},
					},
				},
				{
					Name: "number_list_prop_required",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list of numbers, required."),
						ElementType: ir.ElementType{
							Number: &ir.NumberElement{},
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
