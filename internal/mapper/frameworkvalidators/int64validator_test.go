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

func TestInt64ValidatorOneOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		values   []int64
		expected *schema.CustomValidator
	}{
		"nil": {
			values:   nil,
			expected: nil,
		},
		"empty": {
			values:   []int64{},
			expected: nil,
		},
		"one": {
			values: []int64{1},
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
					},
				},
				SchemaDefinition: "int64validator.OneOf(\n1,\n)",
			},
		},
		"multiple": {
			values: []int64{1, 2},
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
					},
				},
				SchemaDefinition: "int64validator.OneOf(\n1,\n2,\n)",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.Int64ValidatorOneOf(testCase.values)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
