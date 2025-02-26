// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"gopkg.in/yaml.v3"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

func TestBuildStringResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"string attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Description: "hey there! I'm a string type, not sensitive, required.",
					}),
					"string_sensitive_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      "password",
						Description: "hey there! I'm a string type, sensitive",
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_prop",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, not sensitive, required."),
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "string_sensitive_prop",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type, sensitive"),
						Sensitive:                pointer(true),
					},
				},
			},
		},
		"string attributes default": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop_required_default_non_empty"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop_default_empty": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"string"},
						Default: &yaml.Node{Kind: yaml.ScalarNode, Value: ""},
					}),
					"string_prop_default_non_empty": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"string"},
						Default: &yaml.Node{Kind: yaml.ScalarNode, Value: "test value"},
					}),
					"string_prop_required_default_non_empty": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"string"},
						Default: &yaml.Node{Kind: yaml.ScalarNode, Value: "test value"},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_prop_default_empty",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.StringDefault{
							Static: pointer(""),
						},
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "string_prop_default_non_empty",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.StringDefault{
							Static: pointer("test value"),
						},
					},
				},
				&attrmapper.ResourceStringAttribute{
					Name: "string_prop_required_default_non_empty",
					StringAttribute: resource.StringAttribute{
						// Intentionally not required due to default
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.StringDefault{
							Static: pointer("test value"),
						},
					},
				},
			},
		},
		"string attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"string"},
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_prop",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with string element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "string_list_prop_required",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of strings, required."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
		},
		"validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
						Enum: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "one"},
							{Kind: yaml.ScalarNode, Value: "two"},
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_prop",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.StringValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
										},
									},
									SchemaDefinition: "stringvalidator.OneOf(\n\"one\",\n\"two\",\n)",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := oas.OASSchema{Schema: testCase.schema}
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
		expectedAttributes attrmapper.DataSourceAttributes
	}{
		"string attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Description: "hey there! I'm a string type, not sensitive, required.",
					}),
					"string_sensitive_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      "password",
						Description: "hey there! I'm a string type, sensitive",
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_prop",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, not sensitive, required."),
					},
				},
				&attrmapper.DataSourceStringAttribute{
					Name: "string_sensitive_prop",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type, sensitive"),
						Sensitive:                pointer(true),
					},
				},
			},
		},
		"string attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"string"},
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_prop",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with string element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListAttribute{
					Name: "string_list_prop",
					ListAttribute: datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				&attrmapper.DataSourceListAttribute{
					Name: "string_list_prop_required",
					ListAttribute: datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of strings, required."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
		},
		"validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
						Enum: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "one"},
							{Kind: yaml.ScalarNode, Value: "two"},
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_prop",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.StringValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
										},
									},
									SchemaDefinition: "stringvalidator.OneOf(\n\"one\",\n\"two\",\n)",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := oas.OASSchema{Schema: testCase.schema}
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

func TestBuildStringProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ProviderAttributes
	}{
		"string attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Description: "hey there! I'm a string type, not sensitive, required.",
					}),
					"string_sensitive_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"string"},
						Format:      "password",
						Description: "hey there! I'm a string type, sensitive",
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderStringAttribute{
					Name: "string_prop",
					StringAttribute: provider.StringAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a string type, not sensitive, required."),
					},
				},
				&attrmapper.ProviderStringAttribute{
					Name: "string_sensitive_prop",
					StringAttribute: provider.StringAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a string type, sensitive"),
						Sensitive:        pointer(true),
					},
				},
			},
		},
		"string attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"string"},
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderStringAttribute{
					Name: "string_prop",
					StringAttribute: provider.StringAttribute{
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with string element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListAttribute{
					Name: "string_list_prop",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				&attrmapper.ProviderListAttribute{
					Name: "string_list_prop_required",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of strings, required."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
		},
		"validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"string_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
						Enum: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "one"},
							{Kind: yaml.ScalarNode, Value: "two"},
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderStringAttribute{
					Name: "string_prop",
					StringAttribute: provider.StringAttribute{
						OptionalRequired: schema.Required,
						Validators: []schema.StringValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
										},
									},
									SchemaDefinition: "stringvalidator.OneOf(\n\"one\",\n\"two\",\n)",
								},
							},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := oas.OASSchema{Schema: testCase.schema}
			attributes, err := schema.BuildProviderAttributes()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(attributes, testCase.expectedAttributes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetStringValidators(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema   oas.OASSchema
		expected []schema.StringValidator
	}{
		"none": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type: []string{"string"},
				},
			},
			expected: nil,
		},
		"enum": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type: []string{"string"},
					Enum: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "one"},
						{Kind: yaml.ScalarNode, Value: "two"},
					},
				},
			},
			expected: []schema.StringValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
							},
						},
						SchemaDefinition: "stringvalidator.OneOf(\n\"one\",\n\"two\",\n)",
					},
				},
			},
		},
		"maxLength": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:      []string{"string"},
					MaxLength: pointer(int64(123)),
				},
			},
			expected: []schema.StringValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
							},
						},
						SchemaDefinition: "stringvalidator.LengthAtMost(123)",
					},
				},
			},
		},
		"maxLength-and-minLength": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:      []string{"string"},
					MinLength: pointer(int64(123)),
					MaxLength: pointer(int64(456)),
				},
			},
			expected: []schema.StringValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
							},
						},
						SchemaDefinition: "stringvalidator.LengthBetween(123, 456)",
					},
				},
			},
		},
		"minLength": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:      []string{"string"},
					MinLength: pointer(int64(123)),
				},
			},
			expected: []schema.StringValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
							},
						},
						SchemaDefinition: "stringvalidator.LengthAtLeast(123)",
					},
				},
			},
		},
		"pattern": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:    []string{"string"},
					Pattern: "^[a-zA-Z0-9]*$",
				},
			},
			expected: []schema.StringValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "regexp",
							},
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
							},
						},
						SchemaDefinition: "stringvalidator.RegexMatches(regexp.MustCompile(\"^[a-zA-Z0-9]*$\"), \"\")",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetStringValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
