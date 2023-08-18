package merge_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/merge"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMergeResourceAttributes_DescriptionPriority(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		target             []resource.Attribute
		mergeSlices        [][]resource.Attribute
		expectedAttributes []resource.Attribute
	}{
		"primitives": {
			target: []resource.Attribute{
				{
					Name: "bool_attribute",
					Bool: &resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "float64_attribute",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "int64_attribute",
					Int64: &resource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "number_attribute",
					Number: &resource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "string_attribute",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
			},
			mergeSlices: [][]resource.Attribute{
				{
					{
						Name: "bool_attribute",
						Bool: &resource.BoolAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "float64_attribute",
						Float64: &resource.Float64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "int64_attribute",
						Int64: &resource.Int64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "number_attribute",
						Number: &resource.NumberAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "string_attribute",
						String: &resource.StringAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
				},
				{
					{
						Name: "bool_attribute",
						Bool: &resource.BoolAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "float64_attribute",
						Float64: &resource.Float64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "int64_attribute",
						Int64: &resource.Int64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "number_attribute",
						Number: &resource.NumberAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "string_attribute",
						String: &resource.StringAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
				},
				{
					{
						Name: "bool_attribute",
						Bool: &resource.BoolAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "float64_attribute",
						Float64: &resource.Float64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "int64_attribute",
						Int64: &resource.Int64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "number_attribute",
						Number: &resource.NumberAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "string_attribute",
						String: &resource.StringAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
				},
			},
			expectedAttributes: []resource.Attribute{
				{
					Name: "bool_attribute",
					Bool: &resource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "float64_attribute",
					Float64: &resource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "int64_attribute",
					Int64: &resource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "number_attribute",
					Number: &resource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "string_attribute",
					String: &resource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
			},
		},
		"nested attribute types": {
			target: []resource.Attribute{
				{
					Name: "single_nested_attribute",
					SingleNested: &resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						Attributes: resource.Attributes{
							resource.Attribute{
								Name: "string_attribute",
								String: &resource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              nil,
								},
							},
						},
					},
				},
				{
					Name: "list_nested_attribute",
					ListNested: &resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              nil,
									},
								},
							},
						},
					},
				},
				{
					Name: "map_nested_attribute",
					MapNested: &resource.MapNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              nil,
									},
								},
							},
						},
					},
				},
				{
					Name: "set_nested_attribute",
					SetNested: &resource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              nil,
									},
								},
							},
						},
					},
				},
			},
			mergeSlices: [][]resource.Attribute{
				{
					{
						Name: "single_nested_attribute",
						SingleNested: &resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							Attributes: resource.Attributes{
								resource.Attribute{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer(""),
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &resource.ListNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer(""),
										},
									},
								},
							},
						},
					},
					{
						Name: "map_nested_attribute",
						MapNested: &resource.MapNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer(""),
										},
									},
								},
							},
						},
					},
					{
						Name: "set_nested_attribute",
						SetNested: &resource.SetNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer(""),
										},
									},
								},
							},
						},
					},
				},
				{
					{
						Name: "single_nested_attribute",
						SingleNested: &resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							Attributes: resource.Attributes{
								resource.Attribute{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("this one!"),
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &resource.ListNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("this one!"),
										},
									},
								},
							},
						},
					},
					{
						Name: "map_nested_attribute",
						MapNested: &resource.MapNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("this one!"),
										},
									},
								},
							},
						},
					},
					{
						Name: "set_nested_attribute",
						SetNested: &resource.SetNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("this one!"),
										},
									},
								},
							},
						},
					},
				},
				{
					{
						Name: "single_nested_attribute",
						SingleNested: &resource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							Attributes: resource.Attributes{
								resource.Attribute{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("no!"),
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &resource.ListNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("no!"),
										},
									},
								},
							},
						},
					},
					{
						Name: "map_nested_attribute",
						MapNested: &resource.MapNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("no!"),
										},
									},
								},
							},
						},
					},
					{
						Name: "set_nested_attribute",
						SetNested: &resource.SetNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							NestedObject: resource.NestedAttributeObject{
								Attributes: []resource.Attribute{
									{
										Name: "string_attribute",
										String: &resource.StringAttribute{
											ComputedOptionalRequired: schema.Required,
											Description:              pointer("no!"),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttributes: []resource.Attribute{
				{
					Name: "single_nested_attribute",
					SingleNested: &resource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						Attributes: resource.Attributes{
							resource.Attribute{
								Name: "string_attribute",
								String: &resource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("this one!"),
								},
							},
						},
					},
				},
				{
					Name: "list_nested_attribute",
					ListNested: &resource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("this one!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "map_nested_attribute",
					MapNested: &resource.MapNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("this one!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "set_nested_attribute",
					SetNested: &resource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "string_attribute",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("this one!"),
									},
								},
							},
						},
					},
				},
			},
		},
		"collection attribute types": {
			target: []resource.Attribute{
				{
					Name: "list_attribute",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &resource.MapAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &resource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
			mergeSlices: [][]resource.Attribute{
				{
					{
						Name: "list_attribute",
						List: &resource.ListAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &resource.MapAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &resource.SetAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
				},
				{
					{
						Name: "list_attribute",
						List: &resource.ListAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &resource.MapAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &resource.SetAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
				},
				{
					{
						Name: "list_attribute",
						List: &resource.ListAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &resource.MapAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &resource.SetAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
				},
			},
			expectedAttributes: []resource.Attribute{
				{
					Name: "list_attribute",
					List: &resource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &resource.MapAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &resource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
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

			got := merge.MergeResourceAttributes(testCase.target, testCase.mergeSlices...)

			if diff := cmp.Diff(*got, testCase.expectedAttributes); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
