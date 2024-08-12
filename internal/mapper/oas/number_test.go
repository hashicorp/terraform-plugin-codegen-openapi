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
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"gopkg.in/yaml.v3"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

func TestBuildNumberResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"float64 attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_prop_required", "float_float64_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceFloat64Attribute{
					Name: "double_float64_prop",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "double_float64_prop_required",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float_float64_prop",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float_float64_prop_required",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"float64 attributes default": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop_required_default_non_zero"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop_default_non_zero": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"number"},
						Format:  "double",
						Default: &yaml.Node{Kind: yaml.ScalarNode, Value: "123.45"},
					}),
					"float64_prop_default_zero": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"number"},
						Format:  "double",
						Default: &yaml.Node{Kind: yaml.ScalarNode, Value: "0.0"},
					}),
					"float64_prop_required_default_non_zero": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"number"},
						Format:  "double",
						Default: &yaml.Node{Kind: yaml.ScalarNode, Value: "123.45"},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_prop_default_non_zero",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.Float64Default{
							Static: pointer(float64(123.45)),
						},
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_prop_default_zero",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.Float64Default{
							Static: pointer(float64(0.0)),
						},
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_prop_required_default_non_zero",
					Float64Attribute: resource.Float64Attribute{
						// Intentionally not required due to default
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.Float64Default{
							Static: pointer(float64(123.45)),
						},
					},
				},
			},
		},
		"float64 attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"number"},
						Format:     "double",
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_prop",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"float64 attribute validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:   []string{"number"},
						Format: "double",
						Enum: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "1.2"},
							{Kind: yaml.ScalarNode, Value: "2.3"},
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_prop",
					Float64Attribute: resource.Float64Attribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type.",
					}),
					"number_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type, required.",
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceNumberAttribute{
					Name: "number_prop",
					NumberAttribute: resource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a number type."),
					},
				},
				&attrmapper.ResourceNumberAttribute{
					Name: "number_prop_required",
					NumberAttribute: resource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a number type, required."),
					},
				},
			},
		},
		"number attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"number"},
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceNumberAttribute{
					Name: "number_prop",
					NumberAttribute: resource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with float64 element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_list_prop_required", "float_float64_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "double_float64_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "double_float64_list_prop_required",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "float_float64_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "float_float64_list_prop_required",
					ListAttribute: resource.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "number_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of numbers."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "number_list_prop_required",
					ListAttribute: resource.ListAttribute{
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
		expectedAttributes attrmapper.DataSourceAttributes
	}{
		"float64 attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_prop_required", "float_float64_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceFloat64Attribute{
					Name: "double_float64_prop",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				&attrmapper.DataSourceFloat64Attribute{
					Name: "double_float64_prop_required",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float_float64_prop",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float_float64_prop_required",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"float64 attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"number"},
						Format:     "double",
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float64_prop",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"float64 attribute validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:   []string{"number"},
						Format: "double",
						Enum: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "1.2"},
							{Kind: yaml.ScalarNode, Value: "2.3"},
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float64_prop",
					Float64Attribute: datasource.Float64Attribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type.",
					}),
					"number_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type, required.",
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceNumberAttribute{
					Name: "number_prop",
					NumberAttribute: datasource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a number type."),
					},
				},
				&attrmapper.DataSourceNumberAttribute{
					Name: "number_prop_required",
					NumberAttribute: datasource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a number type, required."),
					},
				},
			},
		},
		"number attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"number"},
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceNumberAttribute{
					Name: "number_prop",
					NumberAttribute: datasource.NumberAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with float64 element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_list_prop_required", "float_float64_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListAttribute{
					Name: "double_float64_list_prop",
					ListAttribute: datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.DataSourceListAttribute{
					Name: "double_float64_list_prop_required",
					ListAttribute: datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.DataSourceListAttribute{
					Name: "float_float64_list_prop",
					ListAttribute: datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.DataSourceListAttribute{
					Name: "float_float64_list_prop_required",
					ListAttribute: datasource.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceListAttribute{
					Name: "number_list_prop",
					ListAttribute: datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of numbers."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
						},
					},
				},
				&attrmapper.DataSourceListAttribute{
					Name: "number_list_prop_required",
					ListAttribute: datasource.ListAttribute{
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
		expectedAttributes attrmapper.ProviderAttributes
	}{
		"float64 attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_prop_required", "float_float64_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderFloat64Attribute{
					Name: "double_float64_prop",
					Float64Attribute: provider.Float64Attribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a float64 type, from a double."),
					},
				},
				&attrmapper.ProviderFloat64Attribute{
					Name: "double_float64_prop_required",
					Float64Attribute: provider.Float64Attribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a float64 type, from a double, required."),
					},
				},
				&attrmapper.ProviderFloat64Attribute{
					Name: "float_float64_prop",
					Float64Attribute: provider.Float64Attribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a float64 type, from a float."),
					},
				},
				&attrmapper.ProviderFloat64Attribute{
					Name: "float_float64_prop_required",
					Float64Attribute: provider.Float64Attribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a float64 type, from a float, required."),
					},
				},
			},
		},
		"float64 attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"number"},
						Format:     "double",
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderFloat64Attribute{
					Name: "float64_prop",
					Float64Attribute: provider.Float64Attribute{
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"float64 attribute validators": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"float64_prop"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"float64_prop": base.CreateSchemaProxy(&base.Schema{
						Type:   []string{"number"},
						Format: "double",
						Enum: []*yaml.Node{
							{Kind: yaml.ScalarNode, Value: "1.2"},
							{Kind: yaml.ScalarNode, Value: "2.3"},
						},
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderFloat64Attribute{
					Name: "float64_prop",
					Float64Attribute: provider.Float64Attribute{
						OptionalRequired: schema.Required,
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type.",
					}),
					"number_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"number"},
						Description: "hey there! I'm a number type, required.",
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderNumberAttribute{
					Name: "number_prop",
					NumberAttribute: provider.NumberAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a number type."),
					},
				},
				&attrmapper.ProviderNumberAttribute{
					Name: "number_prop_required",
					NumberAttribute: provider.NumberAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a number type, required."),
					},
				},
			},
		},
		"number attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
					"number_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"number"},
						Deprecated: pointer(true),
					}),
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderNumberAttribute{
					Name: "number_prop",
					NumberAttribute: provider.NumberAttribute{
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with float64 element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"double_float64_list_prop_required", "float_float64_list_prop_required"},
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListAttribute{
					Name: "double_float64_list_prop",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.ProviderListAttribute{
					Name: "double_float64_list_prop_required",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of float64s, required."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.ProviderListAttribute{
					Name: "float_float64_list_prop",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of float64s."),
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
					},
				},
				&attrmapper.ProviderListAttribute{
					Name: "float_float64_list_prop_required",
					ListAttribute: provider.ListAttribute{
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
				Properties: orderedmap.ToOrderedMap(map[string]*base.SchemaProxy{
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
				}),
			},
			expectedAttributes: attrmapper.ProviderAttributes{
				&attrmapper.ProviderListAttribute{
					Name: "number_list_prop",
					ListAttribute: provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of numbers."),
						ElementType: schema.ElementType{
							Number: &schema.NumberType{},
						},
					},
				},
				&attrmapper.ProviderListAttribute{
					Name: "number_list_prop_required",
					ListAttribute: provider.ListAttribute{
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
					Enum: []*yaml.Node{
						{Kind: yaml.ScalarNode, Value: "1.2"},
						{Kind: yaml.ScalarNode, Value: "2.3"},
					},
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
