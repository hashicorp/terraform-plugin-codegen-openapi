// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestResourceListNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceListNestedAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
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
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list nested description"),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
			mergeAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "list_nested_attribute",
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

func TestResourceListNestedAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceListNestedAttribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceListNestedAttribute{
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
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
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new description"),
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyOverride(testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestResourceListNestedAttribute_ApplyNestedOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceListNestedAttribute
		overridePath      []string
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override nested attribute": {
			attribute: attrmapper.ResourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceListNestedAttribute{
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
							ListNestedAttribute: resource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},

				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceListNestedAttribute{
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
							ListNestedAttribute: resource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"override double nested attribute": {
			attribute: attrmapper.ResourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceListNestedAttribute{
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
							ListNestedAttribute: resource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute", "double_nested_attribute"},
			override: explorer.Override{
				Description:              "new description",
				ComputedOptionalRequired: string(schema.Optional),
			},
			expectedAttribute: &attrmapper.ResourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.ResourceNestedAttributeObject{
					Attributes: attrmapper.ResourceAttributes{
						&attrmapper.ResourceListNestedAttribute{
							Name: "nested_attribute",
							NestedObject: attrmapper.ResourceNestedAttributeObject{
								Attributes: attrmapper.ResourceAttributes{
									&attrmapper.ResourceStringAttribute{
										Name: "double_nested_attribute",
										StringAttribute: resource.StringAttribute{
											ComputedOptionalRequired: schema.Optional,
											Description:              pointer("new description"),
										},
									},
								},
							},
							ListNestedAttribute: resource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyNestedOverride(testCase.overridePath, testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDataSourceListNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceListNestedAttribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
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
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list nested description"),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
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
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
			mergeAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "list_nested_attribute",
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

func TestDataSourceListNestedAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceListNestedAttribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceListNestedAttribute{
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
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
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new description"),
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyOverride(testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestDataSourceListNestedAttribute_ApplyNestedOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceListNestedAttribute
		overridePath      []string
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override nested attribute": {
			attribute: attrmapper.DataSourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceListNestedAttribute{
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
							ListNestedAttribute: datasource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},

				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceListNestedAttribute{
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
							ListNestedAttribute: datasource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new description"),
							},
						},
					},
				},
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"override double nested attribute": {
			attribute: attrmapper.DataSourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceListNestedAttribute{
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
							ListNestedAttribute: datasource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			overridePath: []string{"nested_attribute", "double_nested_attribute"},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceListNestedAttribute{
				Name: "attribute",
				NestedObject: attrmapper.DataSourceNestedAttributeObject{
					Attributes: attrmapper.DataSourceAttributes{
						&attrmapper.DataSourceListNestedAttribute{
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
							ListNestedAttribute: datasource.ListNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old description"),
							},
						},
					},
				},
				ListNestedAttribute: datasource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := testCase.attribute.ApplyNestedOverride(testCase.overridePath, testCase.override)

			if diff := cmp.Diff(got, testCase.expectedAttribute); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}
