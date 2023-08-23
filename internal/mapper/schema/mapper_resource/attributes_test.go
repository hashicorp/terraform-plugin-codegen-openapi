// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperAttributes_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttributes     mapper_resource.MapperAttributes
		mergeAttributeSlices []mapper_resource.MapperAttributes
		expectedAttributes   mapper_resource.MapperAttributes
	}{
		"matches and appends": {
			targetAttributes: mapper_resource.MapperAttributes{
				&mapper_resource.MapperStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("string description"),
						Sensitive:                pointer(true),
					},
				},
			},
			mergeAttributeSlices: []mapper_resource.MapperAttributes{
				{
					&mapper_resource.MapperStringAttribute{
						Name: "string_attribute",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Computed,
							Description:              pointer("this will be ignored"),
							Sensitive:                pointer(false),
						},
					},
					&mapper_resource.MapperBoolAttribute{
						Name: "bool_attribute",
						BoolAttribute: resource.BoolAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("bool description"),
						},
					},
				},
				{
					&mapper_resource.MapperFloat64Attribute{
						Name: "float64_attribute",
						Float64Attribute: resource.Float64Attribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("float64 description"),
						},
					},
				},
			},
			expectedAttributes: mapper_resource.MapperAttributes{
				&mapper_resource.MapperStringAttribute{
					Name: "string_attribute",
					StringAttribute: resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("string description"),
						Sensitive:                pointer(true),
					},
				},
				&mapper_resource.MapperBoolAttribute{
					Name: "bool_attribute",
					BoolAttribute: resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("bool description"),
					},
				},
				&mapper_resource.MapperFloat64Attribute{
					Name: "float64_attribute",
					Float64Attribute: resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("float64 description"),
					},
				},
			},
		},
		"recursive - matches and appends": {
			targetAttributes: mapper_resource.MapperAttributes{
				&mapper_resource.MapperSingleNestedAttribute{
					Name: "single_nested_attribute",
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "string_attribute",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("string description"),
								Sensitive:                pointer(true),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("single nested description"),
					},
				},
			},
			mergeAttributeSlices: []mapper_resource.MapperAttributes{
				{
					&mapper_resource.MapperSingleNestedAttribute{
						Name: "single_nested_attribute",
						Attributes: mapper_resource.MapperAttributes{
							&mapper_resource.MapperStringAttribute{
								Name: "string_attribute",
								StringAttribute: resource.StringAttribute{
									ComputedOptionalRequired: schema.Computed,
									Description:              pointer("this will be ignored"),
									Sensitive:                pointer(false),
								},
							},
							&mapper_resource.MapperBoolAttribute{
								Name: "bool_attribute",
								BoolAttribute: resource.BoolAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("bool description"),
								},
							},
						},
						SingleNestedAttribute: resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("single nested description"),
						},
					},
				},
				{
					&mapper_resource.MapperSingleNestedAttribute{
						Name: "single_nested_attribute",
						Attributes: mapper_resource.MapperAttributes{
							&mapper_resource.MapperFloat64Attribute{
								Name: "float64_attribute",
								Float64Attribute: resource.Float64Attribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("float64 description"),
								},
							},
						},
						SingleNestedAttribute: resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("single nested description"),
						},
					},
				},
			},
			expectedAttributes: mapper_resource.MapperAttributes{
				&mapper_resource.MapperSingleNestedAttribute{
					Name: "single_nested_attribute",
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "string_attribute",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("string description"),
								Sensitive:                pointer(true),
							},
						},
						&mapper_resource.MapperBoolAttribute{
							Name: "bool_attribute",
							BoolAttribute: resource.BoolAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("bool description"),
							},
						},
						&mapper_resource.MapperFloat64Attribute{
							Name: "float64_attribute",
							Float64Attribute: resource.Float64Attribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("float64 description"),
							},
						},
					},
					SingleNestedAttribute: resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("single nested description"),
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.targetAttributes.Merge(testCase.mergeAttributeSlices...)

			if diff := cmp.Diff(got, testCase.expectedAttributes); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
