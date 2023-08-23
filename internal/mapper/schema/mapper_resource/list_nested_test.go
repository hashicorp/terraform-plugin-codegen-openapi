package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperListNestedAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   mapper_resource.MapperListNestedAttribute
		mergeAttribute    mapper_resource.MapperAttribute
		expectedAttribute mapper_resource.MapperAttribute
	}{
		"mismatch nested attribute type - no merge": {
			targetAttribute: mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperSetNestedAttribute{
				Name: "set_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
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
			expectedAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list nested description"),
				},
			},
			mergeAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("old nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old list nested description"),
				},
			},
		},
		"nil description - merge": {
			targetAttribute: mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
				},
			},
			mergeAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list nested description"),
				},
			},
		},
		"empty description - merge": {
			targetAttribute: mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer(""),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
				},
			},
			mergeAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.ComputedOptional,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new list nested description"),
				},
			},
			expectedAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
					Attributes: mapper_resource.MapperAttributes{
						&mapper_resource.MapperStringAttribute{
							Name: "nested_string",
							StringAttribute: resource.StringAttribute{
								ComputedOptionalRequired: schema.Required,
								Description:              pointer("new nested string description"),
							},
						},
					},
				},
				ListNestedAttribute: resource.ListNestedAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new list nested description"),
				},
			},
		},
		"nested object - merge": {
			targetAttribute: mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
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
			mergeAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
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
			expectedAttribute: &mapper_resource.MapperListNestedAttribute{
				Name: "list_nested_attribute",
				NestedObject: mapper_resource.MapperNestedAttributeObject{
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
