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

func TestResourceFloat64Attribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceFloat64Attribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("string description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old float64 description"),
				},
			},
			mergeAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new float64 description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old float64 description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new float64 description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new float64 description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new float64 description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new float64 description"),
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

func TestResourceFloat64Attribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceFloat64Attribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceFloat64Attribute{
				Name: "test_attribute",
				Float64Attribute: resource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceFloat64Attribute{
				Name: "test_attribute",
				Float64Attribute: resource.Float64Attribute{
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

func TestDataSourceFloat64Attribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceFloat64Attribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("string description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old float64 description"),
				},
			},
			mergeAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new float64 description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old float64 description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new float64 description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new float64 description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new float64 description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "float64_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new float64 description"),
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

func TestDataSourceFloat64Attribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceFloat64Attribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceFloat64Attribute{
				Name: "test_attribute",
				Float64Attribute: datasource.Float64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceFloat64Attribute{
				Name: "test_attribute",
				Float64Attribute: datasource.Float64Attribute{
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
