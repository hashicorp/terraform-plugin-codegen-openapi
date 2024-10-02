// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestResourceAttributes_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttributes     attrmapper.ResourceAttributes
		mergeAttributeSlices []attrmapper.ResourceAttributes
		expectedAttributes   attrmapper.ResourceAttributes
	}{
		"matches and appends": {
			targetAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("string description"),
						Sensitive:                pointer(true),
					},
				},
			},
			mergeAttributeSlices: []attrmapper.ResourceAttributes{
				{
					&attrmapper.ResourceStringAttribute{
						Name: "string_attribute",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Computed,
							Description:              pointer("this will be ignored"),
							Sensitive:                pointer(false),
						},
					},
					&attrmapper.ResourceBoolAttribute{
						Name: "bool_attribute",
						BoolAttribute: resource.BoolAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("bool description"),
						},
					},
				},
				{
					&attrmapper.ResourceFloat64Attribute{
						Name: "float64_attribute",
						Float64Attribute: resource.Float64Attribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("float64 description"),
						},
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("string description"),
						Sensitive:                pointer(true),
					},
				},
				&attrmapper.ResourceBoolAttribute{
					Name: "bool_attribute",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("bool description"),
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("float64 description"),
					},
				},
			},
		},
		"recursive - matches and appends": {
			targetAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "single_nested_attribute",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "string_attribute",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("string description"),
								Sensitive:                pointer(true),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("single nested description"),
					},
				},
			},
			mergeAttributeSlices: []attrmapper.ResourceAttributes{
				{
					&attrmapper.ResourceSingleNestedAttribute{
						Name: "single_nested_attribute",
						Attributes: attrmapper.ResourceAttributes{
							&attrmapper.ResourceStringAttribute{
								Name: "string_attribute",
								StringAttribute: resource.StringAttribute{
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("this will be ignored"),
									Sensitive:                pointer(false),
								},
							},
							&attrmapper.ResourceBoolAttribute{
								Name: "bool_attribute",
								BoolAttribute: resource.BoolAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("bool description"),
								},
							},
						},
						SingleNestedAttribute: resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("single nested description"),
						},
					},
				},
				{
					&attrmapper.ResourceSingleNestedAttribute{
						Name: "single_nested_attribute",
						Attributes: attrmapper.ResourceAttributes{
							&attrmapper.ResourceFloat64Attribute{
								Name: "float64_attribute",
								Float64Attribute: resource.Float64Attribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("float64 description"),
								},
							},
						},
						SingleNestedAttribute: resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("single nested description"),
						},
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "single_nested_attribute",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "string_attribute",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("string description"),
								Sensitive:                pointer(true),
							},
						},
						&attrmapper.ResourceBoolAttribute{
							Name: "bool_attribute",
							BoolAttribute: resource.BoolAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("bool description"),
							},
						},
						&attrmapper.ResourceFloat64Attribute{
							Name: "float64_attribute",
							Float64Attribute: resource.Float64Attribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("float64 description"),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("single nested description"),
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.targetAttributes.Merge(testCase.mergeAttributeSlices...)

			if diff := cmp.Diff(got, testCase.expectedAttributes); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestResourceAttributes_ApplyOverrides(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		overrides          map[string]explorer.Override
		attributes         attrmapper.ResourceAttributes
		expectedAttributes attrmapper.ResourceAttributes
	}{
		// TODO: this may eventually return an error, but for now just returns without modification
		"no matching overrides": {
			overrides: map[string]explorer.Override{
				"": {
					Description: "new description",
				},
				"attribute_that_doesnt_exist": {
					Description: "new description",
				},
				"string_attribute.attribute_that_doesnt_exist": {
					Description: "new description",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
			},
		},
		"matching overrides": {
			overrides: map[string]explorer.Override{
				"string_attribute": {
					Description: "new string description",
				},
				"float64_attribute": {
					Description: "new float64 description",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("new string description"),
					},
				},
				&attrmapper.ResourceFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("new float64 description"),
					},
				},
			},
		},
		"matching overrides computed": {
			overrides: map[string]explorer.Override{
				"string_attribute": {
					ComputedOptionalRequired: "computed",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer(""),
					},
				},
			},
		},
		"matching overrides optional": {
			overrides: map[string]explorer.Override{
				"string_attribute": {
					ComputedOptionalRequired: "optional",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Optional,
						Description:              pointer(""),
					},
				},
			},
		},
		"matching overrides required": {
			overrides: map[string]explorer.Override{
				"string_attribute": {
					ComputedOptionalRequired: "required",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Computed,
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer(""),
					},
				},
			},
		},
		"matching overrides computed_optional": {
			overrides: map[string]explorer.Override{
				"string_attribute": {
					ComputedOptionalRequired: "computed_optional",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Computed,
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer(""),
					},
				},
			},
		},
		"matching nested overrides": {
			overrides: map[string]explorer.Override{
				"single_nested": {
					Description: "new description",
				},
				"single_nested.list_nested": {
					Description: "new description",
				},
				"single_nested.list_nested.string_attribute": {
					Description: "new description",
				},
			},
			attributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "single_nested",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceListNestedAttribute{
							Name: "list_nested",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "string_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							ListNestedAttribute: resource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Optional,
								Description:              pointer("old description"),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Optional,
						Description:              pointer("old description"),
					},
				},
			},
			expectedAttributes: attrmapper.ResourceAttributes{
				&attrmapper.ResourceSingleNestedAttribute{
					Name: "single_nested",
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceListNestedAttribute{
							Name: "list_nested",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "string_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("new description"),
										},
									},
								},
							},
							ListNestedAttribute: resource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Optional,
								Description:              pointer("new description"),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Optional,
						Description:              pointer("new description"),
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attributes.ApplyOverrides(testCase.overrides)

			if diff := cmp.Diff(got, testCase.expectedAttributes); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
