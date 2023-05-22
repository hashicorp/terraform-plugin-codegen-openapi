package schema_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestBuildBoolResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.ResourceAttribute
	}{
		"boolean attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"list attributes with bool element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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

func TestBuildBoolDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.DataSourceAttribute
	}{
		"boolean attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "bool_prop",
					Bool: &ir.DataSourceBoolAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a bool type."),
					},
				},
				{
					Name: "bool_prop_required",
					Bool: &ir.DataSourceBoolAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a bool type, required."),
					},
				},
			},
		},
		"list attributes with bool element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "bool_list_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a list of bools."),
						ElementType: ir.ElementType{
							Bool: &ir.BoolElement{},
						},
					},
				},
				{
					Name: "bool_list_prop_required",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list of bools, required."),
						ElementType: ir.ElementType{
							Bool: &ir.BoolElement{},
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
