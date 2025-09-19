// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frameworkvalidators_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-spec/code"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"

	"github.com/starburstdata/terraform-plugin-codegen-openapi/internal/mapper/frameworkvalidators"
)

func TestStringValidatorLengthAtLeast(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.LengthAtLeast(123)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.StringValidatorLengthAtLeast(testCase.min)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStringValidatorLengthAtMost(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.LengthAtMost(123)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.StringValidatorLengthAtMost(testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStringValidatorLengthBetween(t *testing.T) {
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
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.LengthBetween(123, 456)",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.StringValidatorLengthBetween(testCase.min, testCase.max)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

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

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.StringValidatorOneOf(testCase.values)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestStringValidatorRegexMatches(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		pattern  string
		message  string
		expected *schema.CustomValidator
	}{
		"empty pattern": {
			pattern: "",
			message: "",
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "regexp",
					},
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.RegexMatches(regexp.MustCompile(\"\"), \"\")",
			},
		},
		"pattern": {
			pattern: "^[a-zA-Z0-9]*$",
			message: "",
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "regexp",
					},
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.RegexMatches(regexp.MustCompile(\"^[a-zA-Z0-9]*$\"), \"\")",
			},
		},
		"message": {
			pattern: "^[a-zA-Z0-9]*$",
			message: "must contain alphanumeric characters",
			expected: &schema.CustomValidator{
				Imports: []code.Import{
					{
						Path: "regexp",
					},
					{
						Path: "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator",
					},
				},
				SchemaDefinition: "stringvalidator.RegexMatches(regexp.MustCompile(\"^[a-zA-Z0-9]*$\"), \"must contain alphanumeric characters\")",
			},
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := frameworkvalidators.StringValidatorRegexMatches(testCase.pattern, testCase.message)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
