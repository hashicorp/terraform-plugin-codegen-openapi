package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperMapNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   mapper_resource.MapperMapNestedAttribute
		mergeAttribute    mapper_resource.MapperAttribute
		expectedAttribute mapper_resource.MapperAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperSetNestedAttribute{
				Name: "set_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested string description"),
							},
						},
					},
				},
				SetNestedAttribute: resource.SetNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("set nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map nested description"),
				},
			},
			mergeAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer(""),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				MapNestedAttribute: resource.MapNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
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
										Description:              pointer("nested string description"),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested object description"),
							},
						},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
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
								&mapper_resource.MapperBoolAttribute{
									Name: "double_nested_bool",
									BoolAttribute: resource.BoolAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested bool description"),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("nested object description"),
							},
						},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapNestedAttribute{
				Name: "map_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: []mapper_resource.MapperAttribute{
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
										Description:              pointer("nested string description"),
									},
								},
								&mapper_resource.MapperBoolAttribute{
									Name: "double_nested_bool",
									BoolAttribute: resource.BoolAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested bool description"),
									},
								},
							},
							SingleNestedAttribute: resource.SingleNestedAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("nested object description"),
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
