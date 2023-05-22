package schema_test

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/mapper/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestBuildStringResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.ResourceAttribute
	}{
		"string attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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
		"list attributes with string element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.ResourceAttribute{
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

func TestBuildStringDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]ir.DataSourceAttribute
	}{
		"string attributes": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "string_prop",
					String: &ir.DataSourceStringAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a string type, not sensitive, required."),
						Sensitive:                pointer(false),
					},
				},
				{
					Name: "string_sensitive_prop",
					String: &ir.DataSourceStringAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a string type, sensitive"),
						Sensitive:                pointer(true),
					},
				},
			},
		},
		"list attributes with string element type": {
			schema: &base.Schema{
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
			expectedAttributes: &[]ir.DataSourceAttribute{
				{
					Name: "string_list_prop",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.ComputedOptional,
						Description:              pointer("hey there! I'm a list of strings."),
						ElementType: ir.ElementType{
							String: &ir.StringElement{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &ir.DataSourceListAttribute{
						ComputedOptionalRequired: ir.Required,
						Description:              pointer("hey there! I'm a list of strings, required."),
						ElementType: ir.ElementType{
							String: &ir.StringElement{},
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
