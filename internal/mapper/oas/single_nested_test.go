// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: add error test for nested objects

func TestBuildSingleNestedResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes mapper_resource.MapperAttributes
	}{
		"single nested attributes": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
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
			expectedAttributes: mapper_resource.MapperAttributes{
				&mapper_resource.MapperSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperSingleNestedAttribute{
							Name: "nested_obj_prop_required",
							Attributes: mapper_resource.MapperAttributes{
								&mapper_resource.MapperFloat64Attribute{
									Name: "nested_float64",
									Float64Attribute: resource.Float64Attribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a nested float64 type."),
									},
								},
								&mapper_resource.MapperInt64Attribute{
									Name: "nested_int64_required",
									Int64Attribute: resource.Int64Attribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("hey there! I'm a nested int64 type, required."),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("hey there! I'm a single nested object type, required."),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a single nested object type."),
					},
				},
			},
		},
		"single nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						Required:   []string{"nested_int64_required"},
						Properties: map[string]*base.SchemaProxy{
							"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
							}),
						},
					}),
				},
			},
			expectedAttributes: mapper_resource.MapperAttributes{
				&mapper_resource.MapperSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperInt64Attribute{
							Name: "nested_int64_required",
							Int64Attribute: resource.Int64Attribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
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

func TestBuildSingleNestedDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]datasource.Attribute
	}{
		"single nested attributes": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
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
					Name: "nested_obj_prop",
					SingleNested: &datasource.SingleNestedAttribute{
						Attributes: []datasource.Attribute{
							{
								Name: "nested_obj_prop_required",
								SingleNested: &datasource.SingleNestedAttribute{
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
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("hey there! I'm a single nested object type, required."),
								},
							},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a single nested object type."),
					},
				},
			},
		},
		"single nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						Required:   []string{"nested_int64_required"},
						Properties: map[string]*base.SchemaProxy{
							"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_obj_prop",
					SingleNested: &datasource.SingleNestedAttribute{
						Attributes: []datasource.Attribute{
							{
								Name: "nested_int64_required",
								Int64: &datasource.Int64Attribute{
									ComputedOptionalRequired: schema.Required,
								},
							},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
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

func TestBuildSingleNestedProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]provider.Attribute
	}{
		"single nested attributes": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
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
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "nested_obj_prop",
					SingleNested: &provider.SingleNestedAttribute{
						Attributes: []provider.Attribute{
							{
								Name: "nested_obj_prop_required",
								SingleNested: &provider.SingleNestedAttribute{
									Attributes: []provider.Attribute{
										{
											Name: "nested_float64",
											Float64: &provider.Float64Attribute{
												OptionalRequired: schema.Optional,
												Description:      pointer("hey there! I'm a nested float64 type."),
											},
										},
										{
											Name: "nested_int64_required",
											Int64: &provider.Int64Attribute{
												OptionalRequired: schema.Required,
												Description:      pointer("hey there! I'm a nested int64 type, required."),
											},
										},
									},
									OptionalRequired: schema.Required,
									Description:      pointer("hey there! I'm a single nested object type, required."),
								},
							},
						},
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a single nested object type."),
					},
				},
			},
		},
		"single nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						Required:   []string{"nested_int64_required"},
						Properties: map[string]*base.SchemaProxy{
							"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "nested_obj_prop",
					SingleNested: &provider.SingleNestedAttribute{
						Attributes: []provider.Attribute{
							{
								Name: "nested_int64_required",
								Int64: &provider.Int64Attribute{
									OptionalRequired: schema.Required,
								},
							},
						},
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
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
