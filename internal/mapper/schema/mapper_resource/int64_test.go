package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperInt64Attribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   mapper_resource.MapperInt64Attribute
		mergeAttribute    mapper_resource.MapperAttribute
		expectedAttribute mapper_resource.MapperAttribute
	}{
		"mismatch type - no merge": {
			targetAttribute: mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
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
			expectedAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old int64 description"),
				},
			},
			mergeAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old int64 description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new int64 description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new int64 description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperInt64Attribute{
				Name: "int64_attribute",
				Int64Attribute: resource.Int64Attribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new int64 description"),
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
