// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: add error tests

func TestBuildMapResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"map nested attribute with props": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map nested type.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_password_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"object"},
										Required:    []string{"nested_int64_required"},
										Description: "hey there! I'm a single nested object type.",
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
									"nested_password_required": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"string"},
										Format:      "password",
										Description: "hey there! I'm a nested string type, required.",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceMapNestedAttribute{
					Name: "nested_map_prop",
					NestedObject: attrmapper.ResourceNestedAttributeObject{
						Attributes: attrmapper.ResourceAttributes{
							&attrmapper.ResourceSingleNestedAttribute{
								Name: "nested_obj_prop",
								Attributes: attrmapper.ResourceAttributes{
									&attrmapper.ResourceFloat64Attribute{
										Name: "nested_float64",
										Float64Attribute: resource.Float64Attribute{
											ComputedOptionalRequired: schema.ComputedOptional,
											Description:              pointer("hey there! I'm a nested float64 type."),
										},
									},
									&attrmapper.ResourceInt64Attribute{
										Name: "nested_int64_required",
										Int64Attribute: resource.Int64Attribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("hey there! I'm a nested int64 type, required."),
										},
									},
								},
								SingleNestedAttribute: resource.SingleNestedAttribute{
									ComputedOptionalRequired: schema.ComputedOptional,
									Description:              pointer("hey there! I'm a single nested object type."),
								},
							},
							&attrmapper.ResourceStringAttribute{
								Name: "nested_password_required",
								StringAttribute: resource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Sensitive:                pointer(true),
									Description:              pointer("hey there! I'm a nested string type, required."),
								},
							},
						},
					},
					MapNestedAttribute: resource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map nested type."),
					},
				},
			},
		},
		"map nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
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
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceMapNestedAttribute{
					Name: "nested_map_prop",
					MapNestedAttribute: resource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
					NestedObject: attrmapper.ResourceNestedAttributeObject{
						Attributes: attrmapper.ResourceAttributes{
							&attrmapper.ResourceInt64Attribute{
								Name: "nested_int64_required",
								Int64Attribute: resource.Int64Attribute{
									ComputedOptionalRequired: schema.Required,
								},
							},
						},
					},
				},
			},
		},
		"map nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_map_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:          []string{"object"},
						MinProperties: pointer(int64(1)),
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
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
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceMapNestedAttribute{
					Name: "nested_map_prop_required",
					NestedObject: attrmapper.ResourceNestedAttributeObject{
						Attributes: attrmapper.ResourceAttributes{
							&attrmapper.ResourceInt64Attribute{
								Name: "nested_int64_required",
								Int64Attribute: resource.Int64Attribute{
									ComputedOptionalRequired: schema.Required,
								},
							},
						},
					},
					MapNestedAttribute: resource.MapNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.MapValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
										},
									},
									SchemaDefinition: "mapvalidator.SizeAtLeast(1)",
								},
							},
						},
					},
				},
			},
		},
		"map attributes with element types": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"map_with_strings_required"},
				Properties: map[string]*base.SchemaProxy{
					"map_with_floats": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with floats.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "float",
							}),
						},
					}),
					"map_with_strings_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with strings, required.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string"},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceMapAttribute{
					Name: "map_with_floats",
					MapAttribute: resource.MapAttribute{
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map type with floats."),
					},
				},
				&attrmapper.ResourceMapAttribute{
					Name: "map_with_strings_required",
					MapAttribute: resource.MapAttribute{
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a map type with strings, required."),
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

func TestBuildMapDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.DataSourceAttributes
	}{

		"map nested attribute with props": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map nested type.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_password_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"object"},
										Required:    []string{"nested_int64_required"},
										Description: "hey there! I'm a single nested object type.",
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
									"nested_password_required": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"string"},
										Format:      "password",
										Description: "hey there! I'm a nested string type, required.",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceMapNestedAttribute{
					Name: "nested_map_prop",
					NestedObject: attrmapper.DataSourceNestedAttributeObject{
						Attributes: attrmapper.DataSourceAttributes{
							&attrmapper.DataSourceSingleNestedAttribute{
								Name: "nested_obj_prop",
								Attributes: attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceFloat64Attribute{
										Name: "nested_float64",
										Float64Attribute: datasource.Float64Attribute{
											ComputedOptionalRequired: schema.ComputedOptional,
											Description:              pointer("hey there! I'm a nested float64 type."),
										},
									},
									&attrmapper.DataSourceInt64Attribute{
										Name: "nested_int64_required",
										Int64Attribute: datasource.Int64Attribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("hey there! I'm a nested int64 type, required."),
										},
									},
								},
								SingleNestedAttribute: datasource.SingleNestedAttribute{
									ComputedOptionalRequired: schema.ComputedOptional,
									Description:              pointer("hey there! I'm a single nested object type."),
								},
							},
							&attrmapper.DataSourceStringAttribute{
								Name: "nested_password_required",
								StringAttribute: datasource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Sensitive:                pointer(true),
									Description:              pointer("hey there! I'm a nested string type, required."),
								},
							},
						},
					},
					MapNestedAttribute: datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map nested type."),
					},
				},
			},
		},
		"map nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
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
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceMapNestedAttribute{
					Name: "nested_map_prop",
					MapNestedAttribute: datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
					NestedObject: attrmapper.DataSourceNestedAttributeObject{
						Attributes: attrmapper.DataSourceAttributes{
							&attrmapper.DataSourceInt64Attribute{
								Name: "nested_int64_required",
								Int64Attribute: datasource.Int64Attribute{
									ComputedOptionalRequired: schema.Required,
								},
							},
						},
					},
				},
			},
		},
		"map nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_map_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:          []string{"object"},
						MinProperties: pointer(int64(1)),
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
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
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceMapNestedAttribute{
					Name: "nested_map_prop_required",
					NestedObject: attrmapper.DataSourceNestedAttributeObject{
						Attributes: attrmapper.DataSourceAttributes{
							&attrmapper.DataSourceInt64Attribute{
								Name: "nested_int64_required",
								Int64Attribute: datasource.Int64Attribute{
									ComputedOptionalRequired: schema.Required,
								},
							},
						},
					},
					MapNestedAttribute: datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Validators: []schema.MapValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
										},
									},
									SchemaDefinition: "mapvalidator.SizeAtLeast(1)",
								},
							},
						},
					},
				},
			},
		},
		"map attributes with element types": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"map_with_strings_required"},
				Properties: map[string]*base.SchemaProxy{
					"map_with_floats": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with floats.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "float",
							}),
						},
					}),
					"map_with_strings_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with strings, required.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string"},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceMapAttribute{
					Name: "map_with_floats",
					MapAttribute: datasource.MapAttribute{
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map type with floats."),
					},
				},
				&attrmapper.DataSourceMapAttribute{
					Name: "map_with_strings_required",
					MapAttribute: datasource.MapAttribute{
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a map type with strings, required."),
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

func TestBuildMapProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ProviderAttributes
	}{

		"map nested attribute with props": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map nested type.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_password_required"},
								Properties: map[string]*base.SchemaProxy{
									"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"object"},
										Required:    []string{"nested_int64_required"},
										Description: "hey there! I'm a single nested object type.",
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
									"nested_password_required": base.CreateSchemaProxy(&base.Schema{
										Type:        []string{"string"},
										Format:      "password",
										Description: "hey there! I'm a nested string type, required.",
									}),
								},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderMapNestedAttribute{
					Name: "nested_map_prop",
					NestedObject: attrmapper.ProviderNestedAttributeObject{
						Attributes: attrmapper.ProviderAttributes{
							&attrmapper.ProviderSingleNestedAttribute{
								Name: "nested_obj_prop",
								Attributes: attrmapper.ProviderAttributes{
									&attrmapper.ProviderFloat64Attribute{
										Name: "nested_float64",
										Float64Attribute: provider.Float64Attribute{
											OptionalRequired: schema.Optional,
											Description:      pointer("hey there! I'm a nested float64 type."),
										},
									},
									&attrmapper.ProviderInt64Attribute{
										Name: "nested_int64_required",
										Int64Attribute: provider.Int64Attribute{
											OptionalRequired: schema.Required,
											Description:      pointer("hey there! I'm a nested int64 type, required."),
										},
									},
								},
								SingleNestedAttribute: provider.SingleNestedAttribute{
									OptionalRequired: schema.Optional,
									Description:      pointer("hey there! I'm a single nested object type."),
								},
							},
							&attrmapper.ProviderStringAttribute{
								Name: "nested_password_required",
								StringAttribute: provider.StringAttribute{
									OptionalRequired: schema.Required,
									Sensitive:        pointer(true),
									Description:      pointer("hey there! I'm a nested string type, required."),
								},
							},
						},
					},
					MapNestedAttribute: provider.MapNestedAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a map nested type."),
					},
				},
			},
		},
		"map nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
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
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderMapNestedAttribute{
					Name: "nested_map_prop",
					MapNestedAttribute: provider.MapNestedAttribute{
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
					},
					NestedObject: attrmapper.ProviderNestedAttributeObject{
						Attributes: attrmapper.ProviderAttributes{
							&attrmapper.ProviderInt64Attribute{
								Name: "nested_int64_required",
								Int64Attribute: provider.Int64Attribute{
									OptionalRequired: schema.Required,
								},
							},
						},
					},
				},
			},
		},
		"map nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_map_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:          []string{"object"},
						MinProperties: pointer(int64(1)),
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
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
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderMapNestedAttribute{
					Name: "nested_map_prop_required",
					NestedObject: attrmapper.ProviderNestedAttributeObject{
						Attributes: attrmapper.ProviderAttributes{
							&attrmapper.ProviderInt64Attribute{
								Name: "nested_int64_required",
								Int64Attribute: provider.Int64Attribute{
									OptionalRequired: schema.Required,
								},
							},
						},
					},
					MapNestedAttribute: provider.MapNestedAttribute{
						OptionalRequired: schema.Required,
						Validators: []schema.MapValidator{
							{
								Custom: &schema.CustomValidator{
									Imports: []code.Import{
										{
											Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
										},
									},
									SchemaDefinition: "mapvalidator.SizeAtLeast(1)",
								},
							},
						},
					},
				},
			},
		},
		"map attributes with element types": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"map_with_strings_required"},
				Properties: map[string]*base.SchemaProxy{
					"map_with_floats": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with floats.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"number"},
								Format: "float",
							}),
						},
					}),
					"map_with_strings_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with strings, required.",
						AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"string"},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderMapAttribute{
					Name: "map_with_floats",
					MapAttribute: provider.MapAttribute{
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a map type with floats."),
					},
				},
				&attrmapper.ProviderMapAttribute{
					Name: "map_with_strings_required",
					MapAttribute: provider.MapAttribute{
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a map type with strings, required."),
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

func TestGetMapValidators(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema   oas.OASSchema
		expected []schema.MapValidator
	}{
		"none": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type: []string{"object"},
				},
			},
			expected: nil,
		},
		"maxProperties": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:          []string{"object"},
					MaxProperties: pointer(int64(123)),
				},
			},
			expected: []schema.MapValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
							},
						},
						SchemaDefinition: "mapvalidator.SizeAtMost(123)",
					},
				},
			},
		},
		"maxProperties-and-minProperties": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:          []string{"object"},
					MinProperties: pointer(int64(123)),
					MaxProperties: pointer(int64(456)),
				},
			},
			expected: []schema.MapValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
							},
						},
						SchemaDefinition: "mapvalidator.SizeBetween(123, 456)",
					},
				},
			},
		},
		"minProperties": {
			schema: oas.OASSchema{
				Schema: &base.Schema{
					Type:          []string{"object"},
					MinProperties: pointer(int64(123)),
				},
			},
			expected: []schema.MapValidator{
				{
					Custom: &schema.CustomValidator{
						Imports: []code.Import{
							{
								Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
							},
						},
						SchemaDefinition: "mapvalidator.SizeAtLeast(123)",
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetMapValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
