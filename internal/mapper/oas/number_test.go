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

func TestBuildNumberResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]resource.Attribute
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "double_float64_prop",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				{
					Name: "double_float64_prop_required",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				{
					Name: "float_float64_prop",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				{
					Name: "float_float64_prop_required",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"float64 attribute validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:   []string{"number"},
						Format: "double",
						Enum:   []any{float64(1.2), float64(2.3)},
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "float64_prop",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.Float64Validator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/float64validator",
										},
									},
									SchemaDefinition: "float64validator.OneOf(\n1.2,\n2.3,\n)",
								},
							},
						},
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "number_prop",
					Number: &resource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a number type."),
					},
				},
				{
					Name: "number_prop_required",
					Number: &resource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "double_float64_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "double_float64_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "float_float64_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "float_float64_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "number_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of numbers."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
						},
					},
				},
				{
					Name: "number_list_prop_required",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of numbers, required."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
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

func TestBuildNumberDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]datasource.Attribute
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "double_float64_prop",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				{
					Name: "double_float64_prop_required",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				{
					Name: "float_float64_prop",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				{
					Name: "float_float64_prop_required",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"float64 attribute validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:   []string{"number"},
						Format: "double",
						Enum:   []any{float64(1.2), float64(2.3)},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "float64_prop",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.Float64Validator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/float64validator",
										},
									},
									SchemaDefinition: "float64validator.OneOf(\n1.2,\n2.3,\n)",
								},
							},
						},
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "number_prop",
					Number: &datasource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a number type."),
					},
				},
				{
					Name: "number_prop_required",
					Number: &datasource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "double_float64_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "double_float64_list_prop_required",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "float_float64_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "float_float64_list_prop_required",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "number_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of numbers."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
						},
					},
				},
				{
					Name: "number_list_prop_required",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of numbers, required."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
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

func TestBuildNumberProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]provider.Attribute
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "double_float64_prop",
					Float64: &provider.Float64Attribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				{
					Name: "double_float64_prop_required",
					Float64: &provider.Float64Attribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				{
					Name: "float_float64_prop",
					Float64: &provider.Float64Attribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				{
					Name: "float_float64_prop_required",
					Float64: &provider.Float64Attribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"float64 attribute validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop"},
				Properties: map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:   []string{"number"},
						Format: "double",
						Enum:   []any{float64(1.2), float64(2.3)},
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "float64_prop",
					Float64: &provider.Float64Attribute{
						OptionalRequired: schema.Required,
						Description:      pointer(""),
						Validators: []schema.Float64Validator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/float64validator",
										},
									},
									SchemaDefinition: "float64validator.OneOf(\n1.2,\n2.3,\n)",
								},
							},
						},
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "number_prop",
					Number: &provider.NumberAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a number type."),
					},
				},
				{
					Name: "number_prop_required",
					Number: &provider.NumberAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a number type, required."),
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "double_float64_list_prop",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "double_float64_list_prop_required",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "float_float64_list_prop",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				{
					Name: "float_float64_list_prop_required",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "number_list_prop",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of numbers."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
						},
					},
				},
				{
					Name: "number_list_prop_required",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of numbers, required."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
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

func TestGetFloatValidators(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema   oas.OASSchema
		expected []schema.Float64Validator
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
					Enum: []any{float64(1.2), float64(2.3)},
				},
			},
			expected: []schema.Float64Validator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/float64validator",
							},
						},
						SchemaDefinition: "float64validator.OneOf(\n1.2,\n2.3,\n)",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetFloatValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
