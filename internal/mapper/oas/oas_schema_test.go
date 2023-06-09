// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: holding here for safe-keeping
func pointer[T any](value T) *T {
	return &value
}

func TestGetDescription_Override(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema              oas.OASSchema
		expectedDescription string
	}{
		"override description": {
			schema: oas.OASSchema{
				SchemaOpts: oas.SchemaOpts{
					OverrideDescription: "this is the correct description!",
				},
				Schema: &base.Schema{
					Description: "this shouldn't show up!",
				},
			},
			expectedDescription: "this is the correct description!",
		},
		"no override of description": {
			schema: oas.OASSchema{
				SchemaOpts: oas.SchemaOpts{},
				Schema: &base.Schema{
					Description: "this is the correct description!",
				},
			},
			expectedDescription: "this is the correct description!",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.schema.GetDescription()
			if *got != testCase.expectedDescription {
				t.Fatalf("unexpected difference, got: %s, wanted: %s", *got, testCase.expectedDescription)
			}
		})
	}
}
