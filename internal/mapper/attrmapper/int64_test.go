// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
)

func TestResourceInt64Attribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceInt64Attribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
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
			expectedAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old int64 description"),
				},
			},
			mergeAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old int64 description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new int64 description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new int64 description"),
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

func TestResourceInt64Attribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceInt64Attribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceInt64Attribute{
				Name: "test_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceInt64Attribute{
				Name: "test_attribute",
				Int64Attribute: resource.Int64Attribute{
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

func TestDataSourceInt64Attribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceInt64Attribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
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
			expectedAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old int64 description"),
				},
			},
			mergeAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old int64 description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new int64 description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new int64 description"),
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

func TestDataSourceInt64Attribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceInt64Attribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceInt64Attribute{
				Name: "test_attribute",
				Int64Attribute: datasource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceInt64Attribute{
				Name: "test_attribute",
				Int64Attribute: datasource.Int64Attribute{
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
