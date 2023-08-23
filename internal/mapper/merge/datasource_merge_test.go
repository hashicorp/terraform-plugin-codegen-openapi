// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package merge_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/merge"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func TestMergeDataSourceAttributes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		target             []datasource.Attribute
		mergeSlices        [][]datasource.Attribute
		expectedAttributes []datasource.Attribute
	}{
		"no matches - appends": {
			target: []datasource.Attribute{},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "string_attribute",
						String: &datasource.StringAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("string!"),
						},
					},
				},
				{
					{
						Name: "bool_attribute",
						Bool: &datasource.BoolAttribute{
							ComputedOptionalRequired: schema.Computed,
							Description:              pointer("bool!"),
						},
					},
				},
				{
					{
						Name: "float64_attribute",
						Float64: &datasource.Float64Attribute{
							ComputedOptionalRequired: schema.Required,
						},
					},
				},
			},
			expectedAttributes: []datasource.Attribute{
				{
					Name: "string_attribute",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("string!"),
					},
				},
				{
					Name: "bool_attribute",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("bool!"),
					},
				},
				{
					Name: "float64_attribute",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
					},
				},
			},
		},
		"nested attributes - recursive appends": {
			target: []datasource.Attribute{
				{
					Name: "map_nested_attribute",
					MapNested: &datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("map nested!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_bool",
									Bool: &datasource.BoolAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested bool!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "list_nested_attribute",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("list nested!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_bool",
									Bool: &datasource.BoolAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested bool!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "set_nested_attribute",
					SetNested: &datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("set nested!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_bool",
									Bool: &datasource.BoolAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested bool!"),
									},
								},
							},
						},
					},
				},
			},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "map_nested_attribute",
						MapNested: &datasource.MapNestedAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("map nested!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "nested_string",
										String: &datasource.StringAttribute{
											ComputedOptionalRequired: schema.ComputedOptional,
											Description:              pointer("nested string!"),
										},
									},
									{
										Name: "nested_bool",
										Bool: &datasource.BoolAttribute{
											ComputedOptionalRequired: schema.Computed,
											Description:              pointer("no!"),
										},
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &datasource.ListNestedAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("list nested!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "nested_string",
										String: &datasource.StringAttribute{
											ComputedOptionalRequired: schema.ComputedOptional,
											Description:              pointer("nested string!"),
										},
									},
									{
										Name: "nested_bool",
										Bool: &datasource.BoolAttribute{
											ComputedOptionalRequired: schema.Computed,
											Description:              pointer("no!"),
										},
									},
								},
							},
						},
					},
					{
						Name: "set_nested_attribute",
						SetNested: &datasource.SetNestedAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("set nested!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "nested_string",
										String: &datasource.StringAttribute{
											ComputedOptionalRequired: schema.ComputedOptional,
											Description:              pointer("nested string!"),
										},
									},
									{
										Name: "nested_bool",
										Bool: &datasource.BoolAttribute{
											ComputedOptionalRequired: schema.Computed,
											Description:              pointer("no!"),
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttributes: []datasource.Attribute{

				{
					Name: "map_nested_attribute",
					MapNested: &datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("map nested!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_bool",
									Bool: &datasource.BoolAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested bool!"),
									},
								},
								{
									Name: "nested_string",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested string!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "list_nested_attribute",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("list nested!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_bool",
									Bool: &datasource.BoolAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested bool!"),
									},
								},
								{
									Name: "nested_string",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested string!"),
									},
								},
							},
						},
					},
				},
				{
					Name: "set_nested_attribute",
					SetNested: &datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("set nested!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_bool",
									Bool: &datasource.BoolAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("nested bool!"),
									},
								},
								{
									Name: "nested_string",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("nested string!"),
									},
								},
							},
						},
					},
				},
			},
		},
		"collection attributes - recursive appends": {
			target: []datasource.Attribute{
				{
					Name: "list_attribute",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("list!"),
						ElementType: schema.ElementType{
							Object: &schema.ObjectType{
								AttributeTypes: []schema.ObjectAttributeType{
									{
										Name: "nested_bool",
										Bool: &schema.BoolType{},
									},
								},
							},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &datasource.MapAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("map!"),
						ElementType: schema.ElementType{
							Object: &schema.ObjectType{
								AttributeTypes: []schema.ObjectAttributeType{
									{
										Name: "nested_bool",
										Bool: &schema.BoolType{},
									},
								},
							},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("set!"),
						ElementType: schema.ElementType{
							Object: &schema.ObjectType{
								AttributeTypes: []schema.ObjectAttributeType{
									{
										Name: "nested_bool",
										Bool: &schema.BoolType{},
									},
								},
							},
						},
					},
				},
			},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "list_attribute",
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.Computed,
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_string",
											String: &schema.StringType{},
										},
									},
								},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Computed,
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_string",
											String: &schema.StringType{},
										},
									},
								},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
							ComputedOptionalRequired: schema.Computed,
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:   "nested_string",
											String: &schema.StringType{},
										},
									},
								},
							},
						},
					},
				},
				{
					{
						Name: "list_attribute",
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.Computed,
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64",
											Int64: &schema.Int64Type{},
										},
										{
											Name: "nested_bool",
											Bool: &schema.BoolType{},
										},
									},
								},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Computed,
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64",
											Int64: &schema.Int64Type{},
										},
										{
											Name: "nested_bool",
											Bool: &schema.BoolType{},
										},
									},
								},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
							ComputedOptionalRequired: schema.Computed,
							ElementType: schema.ElementType{
								Object: &schema.ObjectType{
									AttributeTypes: []schema.ObjectAttributeType{
										{
											Name:  "nested_int64",
											Int64: &schema.Int64Type{},
										},
										{
											Name: "nested_bool",
											Bool: &schema.BoolType{},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttributes: []datasource.Attribute{
				{
					Name: "list_attribute",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("list!"),
						ElementType: schema.ElementType{
							Object: &schema.ObjectType{
								AttributeTypes: []schema.ObjectAttributeType{
									{
										Name: "nested_bool",
										Bool: &schema.BoolType{},
									},
									{
										Name:   "nested_string",
										String: &schema.StringType{},
									},
									{
										Name:  "nested_int64",
										Int64: &schema.Int64Type{},
									},
								},
							},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &datasource.MapAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("map!"),
						ElementType: schema.ElementType{
							Object: &schema.ObjectType{
								AttributeTypes: []schema.ObjectAttributeType{
									{
										Name: "nested_bool",
										Bool: &schema.BoolType{},
									},
									{
										Name:   "nested_string",
										String: &schema.StringType{},
									},
									{
										Name:  "nested_int64",
										Int64: &schema.Int64Type{},
									},
								},
							},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("set!"),
						ElementType: schema.ElementType{
							Object: &schema.ObjectType{
								AttributeTypes: []schema.ObjectAttributeType{
									{
										Name: "nested_bool",
										Bool: &schema.BoolType{},
									},
									{
										Name:   "nested_string",
										String: &schema.StringType{},
									},
									{
										Name:  "nested_int64",
										Int64: &schema.Int64Type{},
									},
								},
							},
						},
					},
				},
			},
		},
		"collection attributes - multi-dimensional recursive appends": {
			target: []datasource.Attribute{
				{
					Name: "list_attribute",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("list!"),
						ElementType: schema.ElementType{
							Map: &schema.MapType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_bool",
												Bool: &schema.BoolType{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &datasource.MapAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("map!"),
						ElementType: schema.ElementType{
							Set: &schema.SetType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_bool",
												Bool: &schema.BoolType{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("set!"),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_bool",
												Bool: &schema.BoolType{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "list_attribute",
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("list!"),
							ElementType: schema.ElementType{
								Map: &schema.MapType{
									ElementType: schema.ElementType{
										Object: &schema.ObjectType{
											AttributeTypes: []schema.ObjectAttributeType{
												{
													Name:   "nested_string",
													String: &schema.StringType{},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Computed,
							Description:              pointer("map!"),
							ElementType: schema.ElementType{
								Set: &schema.SetType{
									ElementType: schema.ElementType{
										Object: &schema.ObjectType{
											AttributeTypes: []schema.ObjectAttributeType{
												{
													Name:   "nested_string",
													String: &schema.StringType{},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("set!"),
							ElementType: schema.ElementType{
								List: &schema.ListType{
									ElementType: schema.ElementType{
										Object: &schema.ObjectType{
											AttributeTypes: []schema.ObjectAttributeType{
												{
													Name:   "nested_string",
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
				{
					{
						Name: "list_attribute",
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.ComputedOptional,
							Description:              pointer("list!"),
							ElementType: schema.ElementType{
								Map: &schema.MapType{
									ElementType: schema.ElementType{
										Object: &schema.ObjectType{
											AttributeTypes: []schema.ObjectAttributeType{
												{
													Name:  "nested_int64",
													Int64: &schema.Int64Type{},
												},
												{
													Name: "nested_bool",
													Bool: &schema.BoolType{},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Computed,
							Description:              pointer("map!"),
							ElementType: schema.ElementType{
								Set: &schema.SetType{
									ElementType: schema.ElementType{
										Object: &schema.ObjectType{
											AttributeTypes: []schema.ObjectAttributeType{
												{
													Name:  "nested_int64",
													Int64: &schema.Int64Type{},
												},
												{
													Name: "nested_bool",
													Bool: &schema.BoolType{},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("set!"),
							ElementType: schema.ElementType{
								List: &schema.ListType{
									ElementType: schema.ElementType{
										Object: &schema.ObjectType{
											AttributeTypes: []schema.ObjectAttributeType{
												{
													Name:  "nested_int64",
													Int64: &schema.Int64Type{},
												},
												{
													Name: "nested_bool",
													Bool: &schema.BoolType{},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedAttributes: []datasource.Attribute{
				{
					Name: "list_attribute",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("list!"),
						ElementType: schema.ElementType{
							Map: &schema.MapType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_bool",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "nested_string",
												String: &schema.StringType{},
											},
											{
												Name:  "nested_int64",
												Int64: &schema.Int64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &datasource.MapAttribute{
						ComputedOptionalRequired: schema.Computed,
						Description:              pointer("map!"),
						ElementType: schema.ElementType{
							Set: &schema.SetType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_bool",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "nested_string",
												String: &schema.StringType{},
											},
											{
												Name:  "nested_int64",
												Int64: &schema.Int64Type{},
											},
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("set!"),
						ElementType: schema.ElementType{
							List: &schema.ListType{
								ElementType: schema.ElementType{
									Object: &schema.ObjectType{
										AttributeTypes: []schema.ObjectAttributeType{
											{
												Name: "nested_bool",
												Bool: &schema.BoolType{},
											},
											{
												Name:   "nested_string",
												String: &schema.StringType{},
											},
											{
												Name:  "nested_int64",
												Int64: &schema.Int64Type{},
											},
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

			got := merge.MergeDataSourceAttributes(testCase.target, testCase.mergeSlices...)

			if diff := cmp.Diff(*got, testCase.expectedAttributes); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func TestMergeDataSourceAttributes_DescriptionPriority(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		target             []datasource.Attribute
		mergeSlices        [][]datasource.Attribute
		expectedAttributes []datasource.Attribute
	}{
		"primitives": {
			target: []datasource.Attribute{
				{
					Name: "bool_attribute",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "float64_attribute",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "int64_attribute",
					Int64: &datasource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "number_attribute",
					Number: &datasource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
				{
					Name: "string_attribute",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
					},
				},
			},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "bool_attribute",
						Bool: &datasource.BoolAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "float64_attribute",
						Float64: &datasource.Float64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "int64_attribute",
						Int64: &datasource.Int64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "number_attribute",
						Number: &datasource.NumberAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
					{
						Name: "string_attribute",
						String: &datasource.StringAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
						},
					},
				},
				{
					{
						Name: "bool_attribute",
						Bool: &datasource.BoolAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "float64_attribute",
						Float64: &datasource.Float64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "int64_attribute",
						Int64: &datasource.Int64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "number_attribute",
						Number: &datasource.NumberAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
					{
						Name: "string_attribute",
						String: &datasource.StringAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
						},
					},
				},
				{
					{
						Name: "bool_attribute",
						Bool: &datasource.BoolAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "float64_attribute",
						Float64: &datasource.Float64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "int64_attribute",
						Int64: &datasource.Int64Attribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "number_attribute",
						Number: &datasource.NumberAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
					{
						Name: "string_attribute",
						String: &datasource.StringAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
						},
					},
				},
			},
			expectedAttributes: []datasource.Attribute{
				{
					Name: "bool_attribute",
					Bool: &datasource.BoolAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "float64_attribute",
					Float64: &datasource.Float64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "int64_attribute",
					Int64: &datasource.Int64Attribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "number_attribute",
					Number: &datasource.NumberAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
				{
					Name: "string_attribute",
					String: &datasource.StringAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
					},
				},
			},
		},
		"nested attribute types": {
			target: []datasource.Attribute{
				{
					Name: "single_nested_attribute",
					SingleNested: &datasource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						Attributes: datasource.Attributes{
							datasource.Attribute{
								Name: "string_attribute",
								String: &datasource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              nil,
								},
							},
						},
					},
				},
				{
					Name: "list_nested_attribute",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
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
					MapNested: &datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
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
					SetNested: &datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              nil,
									},
								},
							},
						},
					},
				},
			},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "single_nested_attribute",
						SingleNested: &datasource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							Attributes: datasource.Attributes{
								datasource.Attribute{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer(""),
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &datasource.ListNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						MapNested: &datasource.MapNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						SetNested: &datasource.SetNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer(""),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						SingleNested: &datasource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							Attributes: datasource.Attributes{
								datasource.Attribute{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("this one!"),
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &datasource.ListNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						MapNested: &datasource.MapNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						SetNested: &datasource.SetNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("this one!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						SingleNested: &datasource.SingleNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							Attributes: datasource.Attributes{
								datasource.Attribute{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Description:              pointer("no!"),
									},
								},
							},
						},
					},
					{
						Name: "list_nested_attribute",
						ListNested: &datasource.ListNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						MapNested: &datasource.MapNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
						SetNested: &datasource.SetNestedAttribute{
							ComputedOptionalRequired: schema.Required,
							Description:              pointer("no!"),
							NestedObject: datasource.NestedAttributeObject{
								Attributes: []datasource.Attribute{
									{
										Name: "string_attribute",
										String: &datasource.StringAttribute{
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
			expectedAttributes: []datasource.Attribute{
				{
					Name: "single_nested_attribute",
					SingleNested: &datasource.SingleNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						Attributes: datasource.Attributes{
							datasource.Attribute{
								Name: "string_attribute",
								String: &datasource.StringAttribute{
									ComputedOptionalRequired: schema.Required,
									Description:              pointer("this one!"),
								},
							},
						},
					},
				},
				{
					Name: "list_nested_attribute",
					ListNested: &datasource.ListNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
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
					MapNested: &datasource.MapNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
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
					SetNested: &datasource.SetNestedAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "string_attribute",
									String: &datasource.StringAttribute{
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
			target: []datasource.Attribute{
				{
					Name: "list_attribute",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &datasource.MapAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &datasource.SetAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              nil,
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
			},
			mergeSlices: [][]datasource.Attribute{
				{
					{
						Name: "list_attribute",
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer(""),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
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
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("this one!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
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
						List: &datasource.ListAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "map_attribute",
						Map: &datasource.MapAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
					{
						Name: "set_attribute",
						Set: &datasource.SetAttribute{
							ComputedOptionalRequired: schema.Optional,
							Description:              pointer("no!"),
							ElementType: schema.ElementType{
								String: &schema.StringType{},
							},
						},
					},
				},
			},
			expectedAttributes: []datasource.Attribute{
				{
					Name: "list_attribute",
					List: &datasource.ListAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "map_attribute",
					Map: &datasource.MapAttribute{
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("this one!"),
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
					},
				},
				{
					Name: "set_attribute",
					Set: &datasource.SetAttribute{
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

			got := merge.MergeDataSourceAttributes(testCase.target, testCase.mergeSlices...)

			if diff := cmp.Diff(*got, testCase.expectedAttributes); diff != "" {
				t.Errorf("Unexpected diagnostics (-got, +expected): %s", diff)
			}
		})
	}
}

func pointer[T any](value T) *T {
	return &value
}
