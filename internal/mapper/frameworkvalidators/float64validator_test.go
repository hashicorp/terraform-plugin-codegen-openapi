// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkvalidators_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
)

func TestFloat64ValidatorOneOf(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		values   []float64
		expected *schema.CustomValidator
	}{
		"nil": {
			values:   nil,
			expected: nil,
		},
		"empty": {
			values:   []float64{},
			expected: nil,
		},
		"one": {
			values: []float64{1.2},
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/float64validator",
					},
				},
				SchemaDefinition: "float64validator.OneOf(\n1.2,\n)",
			},
		},
		"multiple": {
			values: []float64{1.2, 2.3},
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/float64validator",
					},
				},
				SchemaDefinition: "float64validator.OneOf(\n1.2,\n2.3,\n)",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.Float64ValidatorOneOf(testCase.values)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
