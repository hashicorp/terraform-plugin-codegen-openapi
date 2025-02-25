// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.targetAttribute.Merge(testCase.mergeAttribute)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestResourceSetNestedAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceSetNestedAttribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceSetNestedAttribute{
				Name: "test_attribute",
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
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "test_attribute",
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
					Description:              pointer("new description"),
				},
			},
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyOverride(testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestResourceSetNestedAttribute_ApplyNestedOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceSetNestedAttribute
		overridePath      []string
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override nested attribute": {
			attribute: attrmapper.ResourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								Attributes: attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							SetNestedAttribute: resource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},

				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								Attributes: attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							SetNestedAttribute: resource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"override double nested attribute": {
			attribute: attrmapper.ResourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								Attributes: attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							SetNestedAttribute: resource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute", "double_nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								Attributes: attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("new description"),
										},
									},
								},
							},
							SetNestedAttribute: resource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyNestedOverride(testCase.overridePath, testCase.override)

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.targetAttribute.Merge(testCase.mergeAttribute)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDataSourceSetNestedAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceSetNestedAttribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "test_attribute",
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
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "test_attribute",
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
					Description:              pointer("new description"),
				},
			},
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyOverride(testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDataSourceSetNestedAttribute_ApplyNestedOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceSetNestedAttribute
		overridePath      []string
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override nested attribute": {
			attribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.DataSourceNestedAttributeObject{
								Attributes: attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: datasource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							SetNestedAttribute: datasource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},

				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.DataSourceNestedAttributeObject{
								Attributes: attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: datasource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							SetNestedAttribute: datasource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"override double nested attribute": {
			attribute: attrmapper.DataSourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.DataSourceNestedAttributeObject{
								Attributes: attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: datasource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("old description"),
										},
									},
								},
							},
							SetNestedAttribute: datasource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute", "double_nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceSetNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceSetNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.DataSourceNestedAttributeObject{
								Attributes: attrmapper.DataSourceAttributes{
									&attrmapper.DataSourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: datasource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("new description"),
										},
									},
								},
							},
							SetNestedAttribute: datasource.SetNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				SetNestedAttribute: datasource.SetNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
	}
	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyNestedOverride(testCase.overridePath, testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
