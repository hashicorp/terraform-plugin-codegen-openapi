// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/attrmapper"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestResourceNumberAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   attrmapper.ResourceNumberAttribute
		mergeAttribute    attrmapper.ResourceAttribute
		expectedAttribute attrmapper.ResourceAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
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
			expectedAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old number description"),
				},
			},
			mergeAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new number description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old number description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new number description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new number description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new number description"),
				},
			},
			expectedAttribute: &attrmapper.ResourceNumberAttribute{
				Name: "number_attribute",
				NumberAttribute: resource.NumberAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new number description"),
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
