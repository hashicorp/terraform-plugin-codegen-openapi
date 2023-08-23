package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperBoolAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   mapper_resource.MapperBoolAttribute
		mergeAttribute    mapper_resource.MapperAttribute
		expectedAttribute mapper_resource.MapperAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperStringAttribute{
				Name: "string_attribute",
				StringAttribute: resource.StringAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("string description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old bool description"),
				},
			},
			mergeAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old bool description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new bool description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new bool description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperBoolAttribute{
				Name: "bool_attribute",
				BoolAttribute: resource.BoolAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new bool description"),
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
