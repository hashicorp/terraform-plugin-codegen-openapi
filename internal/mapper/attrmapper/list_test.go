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

func TestResourceListAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceListAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch collection type - no merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceSetAttribute{
				Name: "set_attribute",
				SetAttribute: resource.SetAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("set description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"mismatch element type - keep target element type": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Bool: &schema.BoolType{},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name: "nested_bool_attribute",
									Bool: &schema.BoolType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:   "double_nest_string",
												String: &schema.StringType{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name:   "nested_string_attribute",
									String: &schema.StringType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "double_nest_float64",
												Float64: &schema.Float64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name: "nested_bool_attribute",
									Bool: &schema.BoolType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:   "double_nest_string",
												String: &schema.StringType{},
											},
											{
												Name:    "double_nest_float64",
												Float64: &schema.Float64Type{},
											},
										},
									},
								},
								{
									Name:   "nested_string_attribute",
									String: &schema.StringType{},
								},
							},
						},
					},
				},
			},
		},
		"nested list object - merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nested map object - merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nested set object - merge": {
			targetAttribute: attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "list_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
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

func TestResourceListAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceListAttribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceListAttribute{
				Name: "test_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceListAttribute{
				Name: "test_attribute",
				ListAttribute: resource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
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

func TestDataSourceListAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceListAttribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch collection type - no merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceSetAttribute{
				Name: "set_attribute",
				SetAttribute: datasource.SetAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("set description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"mismatch element type - keep target element type": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Bool: &schema.BoolType{},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"nested object - merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name: "nested_bool_attribute",
									Bool: &schema.BoolType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:   "double_nest_string",
												String: &schema.StringType{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name:   "nested_string_attribute",
									String: &schema.StringType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "double_nest_float64",
												Float64: &schema.Float64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name: "nested_bool_attribute",
									Bool: &schema.BoolType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:   "double_nest_string",
												String: &schema.StringType{},
											},
											{
												Name:    "double_nest_float64",
												Float64: &schema.Float64Type{},
											},
										},
									},
								},
								{
									Name:   "nested_string_attribute",
									String: &schema.StringType{},
								},
							},
						},
					},
				},
			},
		},
		"nested list object - merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nested map object - merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nested set object - merge": {
			targetAttribute: attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "list_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
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

func TestDataSourceListAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceListAttribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceListAttribute{
				Name: "test_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceListAttribute{
				Name: "test_attribute",
				ListAttribute: datasource.ListAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
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
