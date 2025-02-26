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

func TestInt64ValidatorAtLeast(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
					},
				},
				SchemaDefinition: "int64validator.AtLeast(123)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.Int64ValidatorAtLeast(testCase.min)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInt64ValidatorAtMost(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
					},
				},
				SchemaDefinition: "int64validator.AtMost(123)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.Int64ValidatorAtMost(testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInt64ValidatorBetween(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/int64validator",
					},
				},
				SchemaDefinition: "int64validator.Between(123, 456)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.Int64ValidatorBetween(testCase.min, testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.Int64ValidatorOneOf(testCase.values)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
