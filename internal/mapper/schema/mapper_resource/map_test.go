package mapper_resource_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/schema/mapper_resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMapperMapAttribute_Merge(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		targetAttribute   mapper_resource.MapperMapAttribute
		mergeAttribute    mapper_resource.MapperAttribute
		expectedAttribute mapper_resource.MapperAttribute
	}{
		"mismatch collection type - no merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperSetAttribute{
				Name: "set_attribute",
				SetAttribute: resource.SetAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("set description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"mismatch element type - keep target element type": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Bool: &schema.BoolType{},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"populated description - no merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("old map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"nil description - merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"empty description - merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer(""),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.ComputedOptional,
					Description:              pointer("new map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					Description:              pointer("new map description"),
					ElementType: schema.ElementType{
						String: &schema.StringType{},
					},
				},
			},
		},
		"nested object - merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name: "nested_bool_attribute",
									Bool: &schema.BoolType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:   "double_nest_string",
												String: &schema.StringType{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name:   "nested_string_attribute",
									String: &schema.StringType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:    "double_nest_float64",
												Float64: &schema.Float64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Object: &schema.ObjectType{
							AttributeTypes: []schema.ObjectAttributeType{
								{
									Name: "nested_bool_attribute",
									Bool: &schema.BoolType{},
								},
								{
									Name: "nested_obj_attribute",
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name:   "double_nest_string",
												String: &schema.StringType{},
											},
											{
												Name:    "double_nest_float64",
												Float64: &schema.Float64Type{},
											},
										},
									},
								},
								{
									Name:   "nested_string_attribute",
									String: &schema.StringType{},
								},
							},
						},
					},
				},
			},
		},
		"nested list object - merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						List: &schema.ListType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nested map object - merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Map: &schema.MapType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"nested set object - merge": {
			targetAttribute: mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttribute: &mapper_resource.MapperMapAttribute{
				Name: "map_attribute",
				MapAttribute: resource.MapAttribute{
					ComputedOptionalRequired: schema.Required,
					ElementType: schema.ElementType{
						Set: &schema.SetType{
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64_attribute",
											Int64: &schema.Int64Type{},
										},
										{
											Name:   "nested_number_attribute",
											Number: &schema.NumberType{},
										},
									},
								},
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
