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

func TestBuildIntegerResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]resource.Attribute
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "int64_prop",
					Int64: &resource.Int64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm an int64 type."),
					},
				},
				{
					Name: "int64_prop_required",
					Int64: &resource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "int64_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of int64s."),
						ElementType: schema.ElementType{
							Int64: &schema.Int64Type{},
						},
					},
				},
				{
					Name: "int64_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of int64s, required."),
						ElementType: schema.ElementType{
							Int64: &schema.Int64Type{},
						},
					},
				},
			},
		},
		"validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"int64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"int64_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"integer"},
						Enum: []any{int64(1), int64(2)},
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "int64_prop",
					Int64: &resource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.Int64Validator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
										},
									},
									SchemaDefinition: "int64validator.OneOf(\n1,\n2,\n)",
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

func TestBuildIntegerDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]datasource.Attribute
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "int64_prop",
					Int64: &datasource.Int64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm an int64 type."),
					},
				},
				{
					Name: "int64_prop_required",
					Int64: &datasource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "int64_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of int64s."),
						ElementType: schema.ElementType{
							Int64: &schema.Int64Type{},
						},
					},
				},
				{
					Name: "int64_list_prop_required",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of int64s, required."),
						ElementType: schema.ElementType{
							Int64: &schema.Int64Type{},
						},
					},
				},
			},
		},
		"validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"int64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"int64_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"integer"},
						Enum: []any{int64(1), int64(2)},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "int64_prop",
					Int64: &datasource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.Int64Validator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
										},
									},
									SchemaDefinition: "int64validator.OneOf(\n1,\n2,\n)",
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

func TestBuildIntegerProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]provider.Attribute
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "int64_prop",
					Int64: &provider.Int64Attribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm an int64 type."),
					},
				},
				{
					Name: "int64_prop_required",
					Int64: &provider.Int64Attribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm an int64 type, required."),
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "int64_list_prop",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of int64s."),
						ElementType: schema.ElementType{
							Int64: &schema.Int64Type{},
						},
					},
				},
				{
					Name: "int64_list_prop_required",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of int64s, required."),
						ElementType: schema.ElementType{
							Int64: &schema.Int64Type{},
						},
					},
				},
			},
		},
		"validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"int64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"int64_prop": base.CreateSchemaProxy(&base.Schema{
						Type: []string{"integer"},
						Enum: []any{int64(1), int64(2)},
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "int64_prop",
					Int64: &provider.Int64Attribute{
						OptionalRequired: schema.Required,
						Validators: []schema.Int64Validator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
										},
									},
									SchemaDefinition: "int64validator.OneOf(\n1,\n2,\n)",
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

func TestGetIntegerValidators(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema   oas.OASSchema
		expected []schema.Int64Validator
	}{
		"none": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type: []string{"integer"},
				},
			},
			expected: nil,
		},
		"enum": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type: []string{"integer"},
					Enum: []any{int64(1), int64(2)},
				},
			},
			expected: []schema.Int64Validator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
							},
						},
						SchemaDefinition: "int64validator.OneOf(\n1,\n2,\n)",
					},
				},
			},
		},
		"maximum": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:    []string{"integer"},
					Maximum: pointer(float64(123)),
				},
			},
			expected: []schema.Int64Validator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
							},
						},
						SchemaDefinition: "int64validator.AtMost(123)",
					},
				},
			},
		},
		"maximum-and-minimum": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:    []string{"integer"},
					Minimum: pointer(float64(123.2)),
					Maximum: pointer(float64(456.2)),
				},
			},
			expected: []schema.Int64Validator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
							},
						},
						SchemaDefinition: "int64validator.Between(123, 456)",
					},
				},
			},
		},
		"minimum": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:    []string{"integer"},
					Minimum: pointer(float64(123)),
				},
			},
			expected: []schema.Int64Validator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
							},
						},
						SchemaDefinition: "int64validator.AtLeast(123)",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetIntegerValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
