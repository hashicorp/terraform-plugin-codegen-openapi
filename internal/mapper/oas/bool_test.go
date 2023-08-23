// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/google/go-cmp/cmp"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

func TestBuildBoolResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes attrmapper.ResourceAttributes
	}{
		"boolean attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type.",
					}),
					"bool_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type, required.",
					}),
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_prop",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a bool type."),
					},
				},
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_prop_required",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a bool type, required."),
					},
				},
			},
		},
		"boolean attributes default": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_prop_required_default_true"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop_default_false": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"boolean"},
						Default: false,
					}),
					"bool_prop_default_true": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"boolean"},
						Default: true,
					}),
					"bool_prop_required_default_true": base.CreateSchemaProxy(&base.Schema{
						Type:    []string{"boolean"},
						Default: true,
					}),
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_prop_default_false",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.BoolDefault{
							Static: pointer(false),
						},
					},
				},
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_prop_default_true",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.BoolDefault{
							Static: pointer(true),
						},
					},
				},
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_prop_required_default_true",
					BoolAttribute: resource.BoolAttribute{
						// Intentionally not required due to default
						ComputedOptionalRequired: schema.ComputedOptional,
						Default: &schema.BoolDefault{
							Static: pointer(true),
						},
					},
				},
			},
		},
		"boolean attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"boolean"},
						Deprecated: pointer(true),
					}),
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_prop",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with bool element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
					"bool_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceListAttribute{
					Name: "bool_list_prop",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of bools."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
						},
					},
				},
				&attrmapper.ResourceListAttribute{
					Name: "bool_list_prop_required",
					ListAttribute: resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of bools, required."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
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

func TestBuildBoolDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]datasource.Attribute
	}{
		"boolean attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type.",
					}),
					"bool_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type, required.",
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "bool_prop",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a bool type."),
					},
				},
				{
					Name: "bool_prop_required",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a bool type, required."),
					},
				},
			},
		},
		"boolean attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"boolean"},
						Deprecated: pointer(true),
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "bool_prop",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						DeprecationMessage:       pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with bool element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
					"bool_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "bool_list_prop",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of bools."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
						},
					},
				},
				{
					Name: "bool_list_prop_required",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a list of bools, required."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
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

func TestBuildBoolProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]provider.Attribute
	}{
		"boolean attributes": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type.",
					}),
					"bool_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"boolean"},
						Description: "hey there! I'm a bool type, required.",
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "bool_prop",
					Bool: &provider.BoolAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a bool type."),
					},
				},
				{
					Name: "bool_prop_required",
					Bool: &provider.BoolAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a bool type, required."),
					},
				},
			},
		},
		"boolean attributes deprecated": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"bool_prop": base.CreateSchemaProxy(&base.Schema{
						Type:       []string{"boolean"},
						Deprecated: pointer(true),
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "bool_prop",
					Bool: &provider.BoolAttribute{
						OptionalRequired:   schema.Optional,
						DeprecationMessage: pointer("This attribute is deprecated."),
					},
				},
			},
		},
		"list attributes with bool element type": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"bool_list_prop_required"},
				Properties: map[string]*base.SchemaProxy{
					"bool_list_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
					"bool_list_prop_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"array"},
						Description: "hey there! I'm a list of bools, required.",
						Items: &base.DynamicValue[*base.SchemaProxy, bool]{
							A: base.CreateSchemaProxy(&base.Schema{
								Type: []string{"boolean"},
							}),
						},
					}),
				},
			},
			expectedAttributes: &[]provider.Attribute{
				{
					Name: "bool_list_prop",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Optional,
						Description:      pointer("hey there! I'm a list of bools."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
						},
					},
				},
				{
					Name: "bool_list_prop_required",
					List: &provider.ListAttribute{
						OptionalRequired: schema.Required,
						Description:      pointer("hey there! I'm a list of bools, required."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
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
