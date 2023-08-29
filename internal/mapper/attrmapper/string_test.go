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

func TestResourceStringAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceStringAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("bool description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old string description"),
				},
			},
			mergeAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new string description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old string description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new string description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new string description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new string description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new string description"),
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

func TestResourceStringAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceStringAttribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceStringAttribute{
				Name: "test_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceStringAttribute{
				Name: "test_attribute",
				StringAttribute: resource.StringAttribute{
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

func TestDataSourceStringAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceStringAttribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("bool description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old string description"),
				},
			},
			mergeAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new string description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old string description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new string description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new string description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new string description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "string_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new string description"),
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

func TestDataSourceStringAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceStringAttribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceStringAttribute{
				Name: "test_attribute",
				StringAttribute: datasource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceStringAttribute{
				Name: "test_attribute",
				StringAttribute: datasource.StringAttribute{
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
