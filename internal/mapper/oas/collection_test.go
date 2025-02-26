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

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// TODO: add error tests

func TestBuildCollectionResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"list nested attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
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
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListNestedAttribute{
					Name: "nested_list_prop_required",
					ListNestedAttribute: resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list nested array type, required."),
					},
					NestedObject: attrmapper.ResourceNestedAttributeObject{
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
					},
				},
			},
		},
		"list nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"array"},
						Deprecated: pointer(true),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListNestedAttribute{
					Name: "nested_list_prop",
					ListNestedAttribute: resource.ListNestedAttribute{
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
		"list nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListNestedAttribute{
					Name: "nested_list_prop_required",
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
					ListNestedAttribute: resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
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
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSetNestedAttribute{
					Name: "nested_set_prop_required",
					SetNestedAttribute: resource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a set nested array type, required."),
					},
					NestedObject: attrmapper.ResourceNestedAttributeObject{
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
					},
				},
			},
		},
		"set nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"array"},
						Format:     "set",
						Deprecated: pointer(true),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSetNestedAttribute{
					Name: "nested_set_prop",
					SetNestedAttribute: resource.SetNestedAttribute{
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
		"set nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						Format:   "set",
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSetNestedAttribute{
					Name: "nested_set_prop_required",
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
					SetNestedAttribute: resource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of lists.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										}),
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "nested_list_prop",
					ListAttribute: resource.ListAttribute{
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
				&attrmapper.ResourceListAttribute{
					Name: "nested_list_prop_required",
					ListAttribute: resource.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										}),
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSetAttribute{
					Name: "nested_set_prop",
					SetAttribute: resource.SetAttribute{
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
				&attrmapper.ResourceSetAttribute{
					Name: "nested_set_prop_required",
					SetAttribute: resource.SetAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"list_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"nested_boolean": base.CreateSchemaProxy(&base.Schema{
												Type:        []string{"boolean"},
												Description: "this won't be added, since it will map to element type",
											}),
											"nested_string": base.CreateSchemaProxy(&base.Schema{
												Type:        []string{"string"},
												Description: "this won't be added, since it will map to element type",
											}),
										}),
									}),
								},
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
								AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "list_with_map",
					ListAttribute: resource.ListAttribute{
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
				&attrmapper.ResourceSetAttribute{
					Name: "set_with_map",
					SetAttribute: resource.SetAttribute{
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
		expectedAttributes attrmapper.DataSourceAttributes
	}{
		"list nested attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
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
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListNestedAttribute{
					Name: "nested_list_prop_required",
					ListNestedAttribute: datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list nested array type, required."),
					},
					NestedObject: attrmapper.DataSourceNestedAttributeObject{
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
					},
				},
			},
		},
		"list nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"array"},
						Deprecated: pointer(true),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListNestedAttribute{
					Name: "nested_list_prop",
					ListNestedAttribute: datasource.ListNestedAttribute{
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
		"list nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListNestedAttribute{
					Name: "nested_list_prop_required",
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
					ListNestedAttribute: datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
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
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSetNestedAttribute{
					Name: "nested_set_prop_required",
					SetNestedAttribute: datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a set nested array type, required."),
					},
					NestedObject: attrmapper.DataSourceNestedAttributeObject{
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
					},
				},
			},
		},
		"set nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"array"},
						Format:     "set",
						Deprecated: pointer(true),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSetNestedAttribute{
					Name: "nested_set_prop",
					SetNestedAttribute: datasource.SetNestedAttribute{
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
		"set nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						Format:   "set",
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSetNestedAttribute{
					Name: "nested_set_prop_required",
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
					SetNestedAttribute: datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of lists.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										}),
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListAttribute{
					Name: "nested_list_prop",
					ListAttribute: datasource.ListAttribute{
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
				&attrmapper.DataSourceListAttribute{
					Name: "nested_list_prop_required",
					ListAttribute: datasource.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										}),
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSetAttribute{
					Name: "nested_set_prop",
					SetAttribute: datasource.SetAttribute{
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
				&attrmapper.DataSourceSetAttribute{
					Name: "nested_set_prop_required",
					SetAttribute: datasource.SetAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"list_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"nested_boolean": base.CreateSchemaProxy(&base.Schema{
												Type:        []string{"boolean"},
												Description: "this won't be added, since it will map to element type",
											}),
											"nested_string": base.CreateSchemaProxy(&base.Schema{
												Type:        []string{"string"},
												Description: "this won't be added, since it will map to element type",
											}),
										}),
									}),
								},
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
								AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListAttribute{
					Name: "list_with_map",
					ListAttribute: datasource.ListAttribute{
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
				&attrmapper.DataSourceSetAttribute{
					Name: "set_with_map",
					SetAttribute: datasource.SetAttribute{
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

func TestBuildCollectionProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ProviderAttributes
	}{
		"list nested attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
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
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListNestedAttribute{
					Name: "nested_list_prop_required",
					ListNestedAttribute: provider.ListNestedAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list nested array type, required."),
					},
					NestedObject: attrmapper.ProviderNestedAttributeObject{
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
					},
				},
			},
		},
		"list nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"array"},
						Deprecated: pointer(true),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListNestedAttribute{
					Name: "nested_list_prop",
					ListNestedAttribute: provider.ListNestedAttribute{
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
		"list nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListNestedAttribute{
					Name: "nested_list_prop_required",
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
					ListNestedAttribute: provider.ListNestedAttribute{
						OptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Format:      "set",
						Description: "hey there! I'm a set nested array type, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
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
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderSetNestedAttribute{
					Name: "nested_set_prop_required",
					SetNestedAttribute: provider.SetNestedAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a set nested array type, required."),
					},
					NestedObject: attrmapper.ProviderNestedAttributeObject{
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
					},
				},
			},
		},
		"set nested attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"array"},
						Format:     "set",
						Deprecated: pointer(true),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderSetNestedAttribute{
					Name: "nested_set_prop",
					SetNestedAttribute: provider.SetNestedAttribute{
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
		"set nested attributes validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"nested_set_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_set_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:     []string{"array"},
						Format:   "set",
						MinItems: pointer(int64(1)),
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type:     []string{"object"},
								Required: []string{"nested_int64_required"},
								Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
									"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
										Type:   []string{"integer"},
										Format: "int64",
									}),
								}),
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderSetNestedAttribute{
					Name: "nested_set_prop_required",
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
					SetNestedAttribute: provider.SetNestedAttribute{
						OptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"nested_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of lists.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"array"},
								Items: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										}),
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListAttribute{
					Name: "nested_list_prop",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of lists."),
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
				&attrmapper.ProviderListAttribute{
					Name: "nested_list_prop_required",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of lists, required."),
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"float64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"number"},
												Format: "double",
											}),
											"int64_prop": base.CreateSchemaProxy(&base.Schema{
												Type:   []string{"integer"},
												Format: "int64",
											}),
										}),
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
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"bool_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"boolean"},
											}),
											"string_prop": base.CreateSchemaProxy(&base.Schema{
												Type: []string{"string"},
											}),
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderSetAttribute{
					Name: "nested_set_prop",
					SetAttribute: provider.SetAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a set of sets."),
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
				&attrmapper.ProviderSetAttribute{
					Name: "nested_set_prop_required",
					SetAttribute: provider.SetAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a set of sets, required."),
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"list_with_map": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list with a nested map of objects.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"object"},
								AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
											"nested_boolean": base.CreateSchemaProxy(&base.Schema{
												Type:        []string{"boolean"},
												Description: "this won't be added, since it will map to element type",
											}),
											"nested_string": base.CreateSchemaProxy(&base.Schema{
												Type:        []string{"string"},
												Description: "this won't be added, since it will map to element type",
											}),
										}),
									}),
								},
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
								AdditionalProperties: &base.DynamicValue[*base.SchemaProxy, bool]{
									A: base.CreateSchemaProxy(&base.Schema{
										Type: []string{"object"},
										Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
										}),
									}),
								},
							}),
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListAttribute{
					Name: "list_with_map",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list with a nested map of objects."),
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
				&attrmapper.ProviderSetAttribute{
					Name: "set_with_map",
					SetAttribute: provider.SetAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a set with a nested map of objects."),
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetSetValidators()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
