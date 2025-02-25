// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// TODO: add error test for nested objects

func TestBuildSingleNestedResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"single nested attributes": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
								}),
							}),
						}),
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceSingleNestedAttribute{
							Name: "nested_obj_prop_required",
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						Required:   []string{"nested_int64_required"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
							}),
						}),
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceInt64Attribute{
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
		expectedAttributes attrmapper.DataSourceAttributes
	}{
		"single nested attributes": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
								}),
							}),
						}),
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceSingleNestedAttribute{
							Name: "nested_obj_prop_required",
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
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("hey there! I'm a single nested object type, required."),
							},
						},
					},
					SingleNestedAttribute: datasource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a single nested object type."),
					},
				},
			},
		},
		"single nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						Required:   []string{"nested_int64_required"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
							}),
						}),
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceInt64Attribute{
							Name: "nested_int64_required",
							Int64Attribute: datasource.Int64Attribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
					SingleNestedAttribute: datasource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
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

func TestBuildSingleNestedProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ProviderAttributes
	}{
		"single nested attributes": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Required:    []string{"nested_obj_prop_required"},
						Description: "hey there! I'm a single nested object type.",
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_obj_prop_required": base.CreateSchemaProxy(&base.Schema{
								Type:        []string{"object"},
								Required:    []string{"nested_int64_required"},
								Description: "hey there! I'm a single nested object type, required.",
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
								}),
							}),
						}),
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: attrmapper.ProviderAttributes{
						&attrmapper.ProviderSingleNestedAttribute{
							Name: "nested_obj_prop_required",
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
								OptionalRequired: schema.Required,
								Description:      pointer("hey there! I'm a single nested object type, required."),
							},
						},
					},
					SingleNestedAttribute: provider.SingleNestedAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a single nested object type."),
					},
				},
			},
		},
		"single nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"object"},
						Deprecated: pointer(true),
						Required:   []string{"nested_int64_required"},
						Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
							"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
								Type:   []string{"integer"},
								Format: "int64",
							}),
						}),
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderSingleNestedAttribute{
					Name: "nested_obj_prop",
					Attributes: attrmapper.ProviderAttributes{
						&attrmapper.ProviderInt64Attribute{
							Name: "nested_int64_required",
							Int64Attribute: provider.Int64Attribute{
								OptionalRequired: schema.Required,
							},
						},
					},
					SingleNestedAttribute: provider.SingleNestedAttribute{
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
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
