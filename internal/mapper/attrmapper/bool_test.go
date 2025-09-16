// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
)

func TestResourceBoolAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceBoolAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
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
			expectedAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old bool description"),
				},
			},
			mergeAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old bool description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new bool description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new bool description"),
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

func TestResourceBoolAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.ResourceBoolAttribute
		override          explorer.Override
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"override description": {
			attribute: attrmapper.ResourceBoolAttribute{
				Name: "test_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.ResourceBoolAttribute{
				Name: "test_attribute",
				BoolAttribute: resource.BoolAttribute{
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

func TestDataSourceBoolAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.DataSourceBoolAttribute
		mergeAttribute    attrmapper.DataSourceAttribute
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
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
			expectedAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old bool description"),
				},
			},
			mergeAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old bool description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new bool description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new bool description"),
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

func TestDataSourceBoolAttribute_ApplyOverride(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		attribute         attrmapper.DataSourceBoolAttribute
		override          explorer.Override
		expectedAttribute attrmapper.DataSourceAttribute
	}{
		"override description": {
			attribute: attrmapper.DataSourceBoolAttribute{
				Name: "test_attribute",
				BoolAttribute: datasource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old description"),
				},
			},
			override: explorer.Override{
				Description: "new description",
			},
			expectedAttribute: &attrmapper.DataSourceBoolAttribute{
				Name: "test_attribute",
				BoolAttribute: datasource.BoolAttribute{
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
