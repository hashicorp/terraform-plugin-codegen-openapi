// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestDataSourceAttributes_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttributes     attrmapper.DataSourceAttributes
		mergeAttributeSlices []attrmapper.DataSourceAttributes
		expectedAttributes   attrmapper.DataSourceAttributes
	}{
		"matches and appends": {
			targetAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("string description"),
						Sensitive:                pointer(true),
					},
				},
			},
			mergeAttributeSlices: []attrmapper.DataSourceAttributes{
				{
					&attrmapper.DataSourceStringAttribute{
						Name: "string_attribute",
						StringAttribute: datasource.StringAttribute{
							ComputedOptionalRequired: schema.Computed,
							Description:              pointer("this will be ignored"),
							Sensitive:                pointer(false),
						},
					},
					&attrmapper.DataSourceBoolAttribute{
						Name: "bool_attribute",
						BoolAttribute: datasource.BoolAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("bool description"),
						},
					},
				},
				{
					&attrmapper.DataSourceFloat64Attribute{
						Name: "float64_attribute",
						Float64Attribute: datasource.Float64Attribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("float64 description"),
						},
					},
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("string description"),
						Sensitive:                pointer(true),
					},
				},
				&attrmapper.DataSourceBoolAttribute{
					Name: "bool_attribute",
					BoolAttribute: datasource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("bool description"),
					},
				},
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("float64 description"),
					},
				},
			},
		},
		"recursive - matches and appends": {
			targetAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSingleNestedAttribute{
					Name: "single_nested_attribute",
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "string_attribute",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("string description"),
								Sensitive:                pointer(true),
							},
						},
					},
					SingleNestedAttribute: datasource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("single nested description"),
					},
				},
			},
			mergeAttributeSlices: []attrmapper.DataSourceAttributes{
				{
					&attrmapper.DataSourceSingleNestedAttribute{
						Name: "single_nested_attribute",
						Attributes: attrmapper.DataSourceAttributes{
							&attrmapper.DataSourceStringAttribute{
								Name: "string_attribute",
								StringAttribute: datasource.StringAttribute{
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("this will be ignored"),
									Sensitive:                pointer(false),
								},
							},
							&attrmapper.DataSourceBoolAttribute{
								Name: "bool_attribute",
								BoolAttribute: datasource.BoolAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("bool description"),
								},
							},
						},
						SingleNestedAttribute: datasource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("single nested description"),
						},
					},
				},
				{
					&attrmapper.DataSourceSingleNestedAttribute{
						Name: "single_nested_attribute",
						Attributes: attrmapper.DataSourceAttributes{
							&attrmapper.DataSourceFloat64Attribute{
								Name: "float64_attribute",
								Float64Attribute: datasource.Float64Attribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("float64 description"),
								},
							},
						},
						SingleNestedAttribute: datasource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("single nested description"),
						},
					},
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSingleNestedAttribute{
					Name: "single_nested_attribute",
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "string_attribute",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("string description"),
								Sensitive:                pointer(true),
							},
						},
						&attrmapper.DataSourceBoolAttribute{
							Name: "bool_attribute",
							BoolAttribute: datasource.BoolAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("bool description"),
							},
						},
						&attrmapper.DataSourceFloat64Attribute{
							Name: "float64_attribute",
							Float64Attribute: datasource.Float64Attribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("float64 description"),
							},
						},
					},
					SingleNestedAttribute: datasource.SingleNestedAttribute{
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

func TestDataSourceAttributes_ApplyOverrides(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		overrides          map[string]explorer.Override
		attributes         attrmapper.DataSourceAttributes
		expectedAttributes attrmapper.DataSourceAttributes
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
			attributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
			},
		},
		"matching overrides": {
			overrides: map[string]explorer.Override{
				"string_attribute": {
					Description:              "new string description",
					ComputedOptionalRequired: "optional",
				},
				"float64_attribute": {
					Description:              "new float64 description",
					ComputedOptionalRequired: "required",
				},
				"computed_optional_attribute": {
					Description:              "new computed_optional",
					ComputedOptionalRequired: "computed_optional",
				},
			},
			attributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
				&attrmapper.DataSourceStringAttribute{
					Name: "computed_optional_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("old description"),
					},
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceStringAttribute{
					Name: "string_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.Optional,
						Description:              pointer("new string description"),
					},
				},
				&attrmapper.DataSourceFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("new float64 description"),
					},
				},
				&attrmapper.DataSourceStringAttribute{
					Name: "computed_optional_attribute",
					StringAttribute: datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("new computed_optional"),
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
			attributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSingleNestedAttribute{
					Name: "single_nested",
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceListNestedAttribute{
							Name: "list_nested",
							NestedObject: attrmapper.DataSourceNestedAttributeObject{
								attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceStringAttribute{
										Name: "string_attribute",
										StringAttribute: datasource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							ListNestedAttribute: datasource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Optional,
								Description:              pointer("old description"),
							},
						},
					},
					SingleNestedAttribute: datasource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Optional,
						Description:              pointer("old description"),
					},
				},
			},
			expectedAttributes: attrmapper.DataSourceAttributes{
				&attrmapper.DataSourceSingleNestedAttribute{
					Name: "single_nested",
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceListNestedAttribute{
							Name: "list_nested",
							NestedObject: attrmapper.DataSourceNestedAttributeObject{
								attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceStringAttribute{
										Name: "string_attribute",
										StringAttribute: datasource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("new description"),
										},
									},
								},
							},
							ListNestedAttribute: datasource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Optional,
								Description:              pointer("new description"),
							},
						},
					},
					SingleNestedAttribute: datasource.SingleNestedAttribute{
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

func pointer[T any](value T) *T {
	return &value
}
