// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: add error tests

func TestBuildCollectionResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]resource.Attribute
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_list_prop_required",
					ListNested: &resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list nested array type, required."),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "nested_float64",
									Float64: &resource.Float64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a nested float64 type."),
									},
								},
								{
									Name: "nested_int64_required",
									Int64: &resource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey there! I'm a nested int64 type, required."),
									},
								},
							},
						},
					},
				},
			},
		},
		"list nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_list_prop_required",
					ListNested: &resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "nested_int64_required",
									Int64: &resource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
									},
								},
							},
						},
						Validators: []schema.ListValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
										},
									},
									SchemaDefinition: "listvalidator.SizeAtLeast(1)",
								},
							},
						},
					},
				},
			},
		},
		"set nested attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set nested array type, required.",
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_set_prop_required",
					SetNested: &resource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a set nested array type, required."),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "nested_float64",
									Float64: &resource.Float64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a nested float64 type."),
									},
								},
								{
									Name: "nested_int64_required",
									Int64: &resource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey there! I'm a nested int64 type, required."),
									},
								},
							},
						},
					},
				},
			},
		},
		"set nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						Format:   "set",
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_set_prop_required",
					SetNested: &resource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "nested_int64_required",
									Int64: &resource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
									},
								},
							},
						},
						Validators: []schema.SetValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator",
										},
									},
									SchemaDefinition: "setvalidator.SizeAtLeast(1)",
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of lists."),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "float64_prop",
												Float64: &schema.Float64Type{},
											},
											{
												Name:  "int64_prop",
												Int64: &schema.Int64Type{},
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
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of lists, required."),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "bool_prop",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "string_prop",
												String: &schema.StringType{},
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
		"set attributes with set and nested object element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_set_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set of sets.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"array"},
								Format: "set",
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
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set of sets, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"array"},
								Format: "set",
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_set_prop",
					Set: &resource.SetAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a set of sets."),
						ElementType: schema.ElementType{
							Set: &schema.SetType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "float64_prop",
												Float64: &schema.Float64Type{},
											},
											{
												Name:  "int64_prop",
												Int64: &schema.Int64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "nested_set_prop_required",
					Set: &resource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a set of sets, required."),
						ElementType: schema.ElementType{
							Set: &schema.SetType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "bool_prop",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "string_prop",
												String: &schema.StringType{},
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
		"list and set attribute - nested map results in element type": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"list_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
									Type: []string{"object"},
									Properties: map[string]*base.SchemaProxy{
										"nested_boolean": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"boolean"},
											Description: "this won't be added, since it will map to element type",
										}),
										"nested_string": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"string"},
											Description: "this won't be added, since it will map to element type",
										}),
									},
								}),
							}),
						},
					}),
					"set_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
									Type: []string{"object"},
									Properties: map[string]*base.SchemaProxy{
										"nested_float64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"number"},
											Format:      "double",
											Description: "this won't be added, since it will map to element type",
										}),
										"nested_int64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"integer"},
											Format:      "int64",
											Description: "this won't be added, since it will map to element type",
										}),
									},
								}),
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "list_with_map",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list with a nested map of objects."),
						ElementType: schema.ElementType{
							Map: &schema.MapType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_boolean",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "nested_string",
												String: &schema.StringType{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "set_with_map",
					Set: &resource.SetAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a set with a nested map of objects."),
						ElementType: schema.ElementType{
							Map: &schema.MapType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "nested_float64",
												Float64: &schema.Float64Type{},
											},
											{
												Name:  "nested_int64",
												Int64: &schema.Int64Type{},
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

func TestBuildCollectionDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]datasource.Attribute
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_list_prop_required",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list nested array type, required."),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_float64",
									Float64: &datasource.Float64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a nested float64 type."),
									},
								},
								{
									Name: "nested_int64_required",
									Int64: &datasource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey there! I'm a nested int64 type, required."),
									},
								},
							},
						},
					},
				},
			},
		},
		"list nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_list_prop_required",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_int64_required",
									Int64: &datasource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
									},
								},
							},
						},
						Validators: []schema.ListValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
										},
									},
									SchemaDefinition: "listvalidator.SizeAtLeast(1)",
								},
							},
						},
					},
				},
			},
		},
		"set nested attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set nested array type, required.",
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_set_prop_required",
					SetNested: &datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a set nested array type, required."),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_float64",
									Float64: &datasource.Float64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a nested float64 type."),
									},
								},
								{
									Name: "nested_int64_required",
									Int64: &datasource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey there! I'm a nested int64 type, required."),
									},
								},
							},
						},
					},
				},
			},
		},
		"set nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						Format:   "set",
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_set_prop_required",
					SetNested: &datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_int64_required",
									Int64: &datasource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
									},
								},
							},
						},
						Validators: []schema.SetValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator",
										},
									},
									SchemaDefinition: "setvalidator.SizeAtLeast(1)",
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of lists."),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "float64_prop",
												Float64: &schema.Float64Type{},
											},
											{
												Name:  "int64_prop",
												Int64: &schema.Int64Type{},
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
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of lists, required."),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "bool_prop",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "string_prop",
												String: &schema.StringType{},
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
		"set attributes with set and nested object element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_set_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set of sets.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"array"},
								Format: "set",
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
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set of sets, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"array"},
								Format: "set",
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
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_set_prop",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a set of sets."),
						ElementType: schema.ElementType{
							Set: &schema.SetType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "float64_prop",
												Float64: &schema.Float64Type{},
											},
											{
												Name:  "int64_prop",
												Int64: &schema.Int64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "nested_set_prop_required",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a set of sets, required."),
						ElementType: schema.ElementType{
							Set: &schema.SetType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "bool_prop",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "string_prop",
												String: &schema.StringType{},
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
		"list and set attribute - nested map results in element type": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"list_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
									Type: []string{"object"},
									Properties: map[string]*base.SchemaProxy{
										"nested_boolean": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"boolean"},
											Description: "this won't be added, since it will map to element type",
										}),
										"nested_string": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"string"},
											Description: "this won't be added, since it will map to element type",
										}),
									},
								}),
							}),
						},
					}),
					"set_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
									Type: []string{"object"},
									Properties: map[string]*base.SchemaProxy{
										"nested_float64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"number"},
											Format:      "double",
											Description: "this won't be added, since it will map to element type",
										}),
										"nested_int64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"integer"},
											Format:      "int64",
											Description: "this won't be added, since it will map to element type",
										}),
									},
								}),
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "list_with_map",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list with a nested map of objects."),
						ElementType: schema.ElementType{
							Map: &schema.MapType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_boolean",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "nested_string",
												String: &schema.StringType{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "set_with_map",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a set with a nested map of objects."),
						ElementType: schema.ElementType{
							Map: &schema.MapType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "nested_float64",
												Float64: &schema.Float64Type{},
											},
											{
												Name:  "nested_int64",
												Int64: &schema.Int64Type{},
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

func TestGetListValidators(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema   oas.OASSchema
		expected []schema.ListValidator
	}{
		"none": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type: []string{"array"},
				},
			},
			expected: nil,
		},
		"maxItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:     []string{"array"},
					MaxItems: pointer(int64(123)),
				},
			},
			expected: []schema.ListValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
							},
						},
						SchemaDefinition: "listvalidator.SizeAtMost(123)",
					},
				},
			},
		},
		"maxItems-and-minItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:     []string{"array"},
					MinItems: pointer(int64(123)),
					MaxItems: pointer(int64(456)),
				},
			},
			expected: []schema.ListValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
							},
						},
						SchemaDefinition: "listvalidator.SizeBetween(123, 456)",
					},
				},
			},
		},
		"minItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:     []string{"array"},
					MinItems: pointer(int64(123)),
				},
			},
			expected: []schema.ListValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
							},
						},
						SchemaDefinition: "listvalidator.SizeAtLeast(123)",
					},
				},
			},
		},
		"uniqueItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:        []string{"array"},
					UniqueItems: pointer(true),
				},
			},
			expected: []schema.ListValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
							},
						},
						SchemaDefinition: "listvalidator.UniqueValues()",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetListValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestGetSetValidators(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema   oas.OASSchema
		expected []schema.SetValidator
	}{
		"none": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:   []string{"array"},
					Format: "set",
				},
			},
			expected: nil,
		},
		"maxItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:     []string{"array"},
					Format:   "set",
					MaxItems: pointer(int64(123)),
				},
			},
			expected: []schema.SetValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator",
							},
						},
						SchemaDefinition: "setvalidator.SizeAtMost(123)",
					},
				},
			},
		},
		"maxItems-and-minItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:     []string{"array"},
					Format:   "set",
					MinItems: pointer(int64(123)),
					MaxItems: pointer(int64(456)),
				},
			},
			expected: []schema.SetValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator",
							},
						},
						SchemaDefinition: "setvalidator.SizeBetween(123, 456)",
					},
				},
			},
		},
		"minItems": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:     []string{"array"},
					Format:   "set",
					MinItems: pointer(int64(123)),
				},
			},
			expected: []schema.SetValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/setvalidator",
							},
						},
						SchemaDefinition: "setvalidator.SizeAtLeast(123)",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetSetValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
