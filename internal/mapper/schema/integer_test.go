package schema_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/ir"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestBuildIntegerResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.ResourceAttribute
	}{
		"int64 attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"list attributes with int64 element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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

func TestBuildIntegerDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.DataSourceAttribute
	}{
		"int64 attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "int64_prop",
					Int64: &ir.DataSourceInt64Attribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm an int64 type."),
					},
				},
				{
					Name: "int64_prop_required",
					Int64: &ir.DataSourceInt64Attribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm an int64 type, required."),
					},
				},
			},
		},
		"list attributes with int64 element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "int64_list_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a list of int64s."),
						ElementType: ir.ElementType{
							Int64: &ir.Int64Element{},
						},
					},
				},
				{
					Name: "int64_list_prop_required",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list of int64s, required."),
						ElementType: ir.ElementType{
							Int64: &ir.Int64Element{},
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
