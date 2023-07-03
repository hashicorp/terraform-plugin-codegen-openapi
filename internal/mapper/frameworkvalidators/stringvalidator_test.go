// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkvalidators_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestStringValidatorOneOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		values   []string
		expected *schema.CustomValidator
	}{
		"nil": {
			values:   nil,
			expected: nil,
		},
		"empty": {
			values:   []string{},
			expected: nil,
		},
		"one": {
			values: []string{"one"},
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.OneOf(\n\"one\",\n)",
			},
		},
		"multiple": {
			values: []string{"one", "two"},
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.OneOf(\n\"one\",\n\"two\",\n)",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.StringValidatorOneOf(testCase.values)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
