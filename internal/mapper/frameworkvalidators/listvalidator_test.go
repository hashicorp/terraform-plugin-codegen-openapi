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

func TestListValidatorSizeAtLeast(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
					},
				},
				SchemaDefinition: "listvalidator.SizeAtLeast(123)",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.ListValidatorSizeAtLeast(testCase.min)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListValidatorSizeAtMost(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
					},
				},
				SchemaDefinition: "listvalidator.SizeAtMost(123)",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.ListValidatorSizeAtMost(testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListValidatorSizeBetween(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
					},
				},
				SchemaDefinition: "listvalidator.SizeBetween(123, 456)",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.ListValidatorSizeBetween(testCase.min, testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestListValidatorUniqueValues(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		expected *schema.CustomValidator
	}{
		"test": {
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator",
					},
				},
				SchemaDefinition: "listvalidator.UniqueValues()",
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.ListValidatorUniqueValues()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
