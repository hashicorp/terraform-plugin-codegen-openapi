// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkvalidators_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
)

func TestMapValidatorSizeAtLeast(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		min      int64
		expected *schema.CustomValidator
	}{
		"test": {
			min: 123,
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
					},
				},
				SchemaDefinition: "mapvalidator.SizeAtLeast(123)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.MapValidatorSizeAtLeast(testCase.min)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMapValidatorSizeAtMost(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		max      int64
		expected *schema.CustomValidator
	}{
		"test": {
			max: 123,
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
					},
				},
				SchemaDefinition: "mapvalidator.SizeAtMost(123)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.MapValidatorSizeAtMost(testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestMapValidatorSizeBetween(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		min      int64
		max      int64
		expected *schema.CustomValidator
	}{
		"test": {
			min: 123,
			max: 456,
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator",
					},
				},
				SchemaDefinition: "mapvalidator.SizeBetween(123, 456)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.MapValidatorSizeBetween(testCase.min, testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
