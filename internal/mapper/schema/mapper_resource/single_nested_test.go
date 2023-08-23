package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperSingleNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   mapper_resource.MapperSingleNestedAttribute
		mergeAttribute    mapper_resource.MapperAttribute
		expectedAttribute mapper_resource.MapperAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("nested string description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("old nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old single nested description"),
				},
			},
			mergeAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("new nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new single nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("old nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old single nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("new nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new single nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("new nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new single nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("new nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new single nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("new nested string description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new single nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("nested string description"),
						},
					},
					&mapper_resource.MapperSingleNestedAttribute{
						Name: "nested_object",
						Attributes: mapper_resource.MapperAttributes{
							&mapper_resource.MapperStringAttribute{
								Name: "double_nested_string",
								StringAttribute: resource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
								},
							},
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperBoolAttribute{
						Name: "nested_bool",
						BoolAttribute: resource.BoolAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("nested bool description"),
						},
					},
					&mapper_resource.MapperSingleNestedAttribute{
						Name: "nested_object",
						Attributes: mapper_resource.MapperAttributes{
							&mapper_resource.MapperStringAttribute{
								Name: "double_nested_string",
								StringAttribute: resource.StringAttribute{
									ComputedOptionalRequired: schema.ComputedOptional,
									Description:              pointer("double nested string description"),
								},
							},
							&mapper_resource.MapperBoolAttribute{
								Name: "double_nested_bool",
								BoolAttribute: resource.BoolAttribute{
									ComputedOptionalRequired: schema.ComputedOptional,
									Description:              pointer("double nested bool description"),
								},
							},
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("single nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperSingleNestedAttribute{
				Name: "single_nested_attribute",
				Attributes: mapper_resource.MapperAttributes{
					&mapper_resource.MapperStringAttribute{
						Name: "nested_string",
						StringAttribute: resource.StringAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("nested string description"),
						},
					},
					&mapper_resource.MapperSingleNestedAttribute{
						Name: "nested_object",
						Attributes: mapper_resource.MapperAttributes{
							&mapper_resource.MapperStringAttribute{
								Name: "double_nested_string",
								StringAttribute: resource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("double nested string description"),
								},
							},
							&mapper_resource.MapperBoolAttribute{
								Name: "double_nested_bool",
								BoolAttribute: resource.BoolAttribute{
									ComputedOptionalRequired: schema.ComputedOptional,
									Description:              pointer("double nested bool description"),
								},
							},
						},
					},
					&mapper_resource.MapperBoolAttribute{
						Name: "nested_bool",
						BoolAttribute: resource.BoolAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("nested bool description"),
						},
					},
				},
				SingleNestedAttribute: resource.SingleNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("single nested description"),
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
