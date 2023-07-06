// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

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
		expectedAttributes *[]resource.Attribute
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "bool_prop",
					Bool: &resource.BoolAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a bool type."),
					},
				},
				{
					Name: "bool_prop_required",
					Bool: &resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a bool type, required."),
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
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "bool_list_prop",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a list of bools."),
						ElementType: schema.ElementType{
							Bool: &schema.BoolType{},
						},
					},
				},
				{
					Name: "bool_list_prop_required",
					List: &resource.ListAttribute{
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
