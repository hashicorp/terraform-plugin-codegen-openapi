// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

func pointer[T any](value T) *T {
	return &value
}
