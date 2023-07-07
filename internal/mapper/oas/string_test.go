// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestBuildStringResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]resource.Attribute
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_prop",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, not sensitive, required."),
					},
				},
				{
					Name: "string_sensitive_prop",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a string type, sensitive"),
						Sensitive:                pointer(true),
					},
				},
			},
		},
		"string attributes default": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"string_prop_default_empty": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"string"},
						Format:  "double",
						Default: "",
					}),
					"string_prop_default_non_empty": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"string"},
						Format:  "double",
						Default: "test value",
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_prop_default_empty",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.StringDefault{
							Static: pointer(""),
						},
					},
				},
				{
					Name: "string_prop_default_non_empty",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.StringDefault{
							Static: pointer("test value"),
						},
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &resource.ListAttribute{
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
				Properties: map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
						Enum: []any{"one", "two"},
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "string_prop",
					String: &resource.StringAttribute{
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
		name, testCase := name, testCase

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
		expectedAttributes *[]datasource.Attribute
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "string_prop",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a string type, not sensitive, required."),
					},
				},
				{
					Name: "string_sensitive_prop",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "string_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &datasource.ListAttribute{
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
				Properties: map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
						Enum: []any{"one", "two"},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "string_prop",
					String: &datasource.StringAttribute{
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
		name, testCase := name, testCase

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
		expectedAttributes *[]provider.Attribute
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "string_prop",
					String: &provider.StringAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a string type, not sensitive, required."),
					},
				},
				{
					Name: "string_sensitive_prop",
					String: &provider.StringAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a string type, sensitive"),
						Sensitive:        pointer(true),
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "string_list_prop",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of strings."),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "string_list_prop_required",
					List: &provider.ListAttribute{
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
				Properties: map[string]*base.SchemaProxy{
					"string_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"string"},
						Enum: []any{"one", "two"},
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "string_prop",
					String: &provider.StringAttribute{
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
		name, testCase := name, testCase

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
					Enum: []any{"one", "two"},
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
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetStringValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
