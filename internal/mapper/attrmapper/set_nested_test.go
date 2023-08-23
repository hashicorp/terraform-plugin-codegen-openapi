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

func TestResourceSetNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceSetNestedAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("list nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old set nested description"),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new set nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old set nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new set nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new set nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new set nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new set nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
			mergeAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "set_nested_attribute",
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

			got, _ := testCase.targetAttribute.Merge(testCase.mergeAttribute)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDataSourceSetNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceSetNestedAttribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("list nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old set nested description"),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new set nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old set nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new set nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new set nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new set nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new set nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
			mergeAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "set_nested_attribute",
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

			got, _ := testCase.targetAttribute.Merge(testCase.mergeAttribute)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
