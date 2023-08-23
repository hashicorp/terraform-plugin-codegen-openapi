// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestResourceMapNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceMapNestedAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("set nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map nested description"),
				},
			},
			mergeAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer(""),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested string description"),
							},
						},
						&attrmapper.ResourceSingleNestedAttribute{
							Name: "nested_object",
							Attributes: attrmapper.ResourceAttributes{
								&attrmapper.ResourceStringAttribute{
									Name: "double_nested_string",
									StringAttribute: resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested string description"),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested object description"),
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceBoolAttribute{
							Name: "nested_bool",
							BoolAttribute: resource.BoolAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested bool description"),
							},
						},
						&attrmapper.ResourceSingleNestedAttribute{
							Name: "nested_object",
							Attributes: attrmapper.ResourceAttributes{
								&attrmapper.ResourceBoolAttribute{
									Name: "double_nested_bool",
									BoolAttribute: resource.BoolAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested bool description"),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested object description"),
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested string description"),
							},
						},
						&attrmapper.ResourceSingleNestedAttribute{
							Name: "nested_object",
							Attributes: attrmapper.ResourceAttributes{
								&attrmapper.ResourceStringAttribute{
									Name: "double_nested_string",
									StringAttribute: resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested string description"),
									},
								},
								&attrmapper.ResourceBoolAttribute{
									Name: "double_nested_bool",
									BoolAttribute: resource.BoolAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested bool description"),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested object description"),
							},
						},
						&attrmapper.ResourceBoolAttribute{
							Name: "nested_bool",
							BoolAttribute: resource.BoolAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested bool description"),
							},
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

			got := testCase.targetAttribute.Merge(testCase.mergeAttribute)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDataSourceMapNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceMapNestedAttribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("set nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map nested description"),
				},
			},
			mergeAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer(""),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: datasource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested string description"),
							},
						},
						&attrmapper.DataSourceSingleNestedAttribute{
							Name: "nested_object",
							Attributes: attrmapper.DataSourceAttributes{
								&attrmapper.DataSourceStringAttribute{
									Name: "double_nested_string",
									StringAttribute: datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested string description"),
									},
								},
							},
							SingleNestedAttribute: datasource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested object description"),
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceBoolAttribute{
							Name: "nested_bool",
							BoolAttribute: datasource.BoolAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested bool description"),
							},
						},
						&attrmapper.DataSourceSingleNestedAttribute{
							Name: "nested_object",
							Attributes: attrmapper.DataSourceAttributes{
								&attrmapper.DataSourceBoolAttribute{
									Name: "double_nested_bool",
									BoolAttribute: datasource.BoolAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested bool description"),
									},
								},
							},
							SingleNestedAttribute: datasource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested object description"),
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceStringAttribute{
							Name: "nested_string",
							StringAttribute: datasource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested string description"),
							},
						},
						&attrmapper.DataSourceSingleNestedAttribute{
							Name: "nested_object",
							Attributes: attrmapper.DataSourceAttributes{
								&attrmapper.DataSourceStringAttribute{
									Name: "double_nested_string",
									StringAttribute: datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested string description"),
									},
								},
								&attrmapper.DataSourceBoolAttribute{
									Name: "double_nested_bool",
									BoolAttribute: datasource.BoolAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested bool description"),
									},
								},
							},
							SingleNestedAttribute: datasource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested object description"),
							},
						},
						&attrmapper.DataSourceBoolAttribute{
							Name: "nested_bool",
							BoolAttribute: datasource.BoolAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested bool description"),
							},
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

			got := testCase.targetAttribute.Merge(testCase.mergeAttribute)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
