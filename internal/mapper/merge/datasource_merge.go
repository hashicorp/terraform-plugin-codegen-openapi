// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package merge

import "github.com/hashicorp/terraform-plugin-codegen-spec/datasource"

// MergeDataSourceAttributes takes a 'targetSlice' of data source attributes, which all other 'mergeSlices' of attributes
// are merged to. Attributes with the same name will be merged together and certain properties, like 'description', will be
// populated based on priority (targetSlice having the highest priority, then mergeSlices[0], mergeSlices[1], etc).
//
// All attributes not present in targetSlice will be appended to the end.
func MergeDataSourceAttributes(targetSlice []datasource.Attribute, mergeSlices ...[]datasource.Attribute) *[]datasource.Attribute {
	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.Name == mergeAttribute.Name {
					// Handle primitive attribute type
					// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#primitive-attribute-types
					if targetAttribute.Bool != nil && mergeAttribute.Bool != nil {
						targetSlice[i].Bool = mergeDataSourceBoolAttribute(*targetAttribute.Bool, *mergeAttribute.Bool)
					} else if targetAttribute.Float64 != nil && mergeAttribute.Float64 != nil {
						targetSlice[i].Float64 = mergeDataSourceFloat64Attribute(*targetAttribute.Float64, *mergeAttribute.Float64)
					} else if targetAttribute.Int64 != nil && mergeAttribute.Int64 != nil {
						targetSlice[i].Int64 = mergeDataSourceInt64Attribute(*targetAttribute.Int64, *mergeAttribute.Int64)
					} else if targetAttribute.Number != nil && mergeAttribute.Number != nil {
						targetSlice[i].Number = mergeDataSourceNumberAttribute(*targetAttribute.Number, *mergeAttribute.Number)
					} else if targetAttribute.String != nil && mergeAttribute.String != nil {
						targetSlice[i].String = mergeDataSourceStringAttribute(*targetAttribute.String, *mergeAttribute.String)
					}

					// Handle nested attribute type
					// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#nested-attribute-types
					if targetAttribute.SingleNested != nil && mergeAttribute.SingleNested != nil {
						targetSlice[i].SingleNested = mergeDataSourceSingleNestedAttribute(*targetAttribute.SingleNested, *mergeAttribute.SingleNested)
					} else if targetAttribute.ListNested != nil && mergeAttribute.ListNested != nil {
						targetSlice[i].ListNested = mergeDataSourceListNestedAttribute(*targetAttribute.ListNested, *mergeAttribute.ListNested)
					} else if targetAttribute.MapNested != nil && mergeAttribute.MapNested != nil {
						targetSlice[i].MapNested = mergeDataSourceMapNestedAttribute(*targetAttribute.MapNested, *mergeAttribute.MapNested)
					} else if targetAttribute.SetNested != nil && mergeAttribute.SetNested != nil {
						targetSlice[i].SetNested = mergeDataSourceSetNestedAttribute(*targetAttribute.SetNested, *mergeAttribute.SetNested)
					}

					// Handle collection attribute type
					// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#collection-attribute-types
					if targetAttribute.List != nil && mergeAttribute.List != nil {
						targetSlice[i].List = mergeDataSourceListAttribute(*targetAttribute.List, *mergeAttribute.List)
					} else if targetAttribute.Map != nil && mergeAttribute.Map != nil {
						targetSlice[i].Map = mergeDataSourceMapAttribute(*targetAttribute.Map, *mergeAttribute.Map)
					} else if targetAttribute.Set != nil && mergeAttribute.Set != nil {
						targetSlice[i].Set = mergeDataSourceSetAttribute(*targetAttribute.Set, *mergeAttribute.Set)
					}

					isNewAttribute = false
					break
				}
			}

			if isNewAttribute {
				// Add this back to the original slice to avoid adding duplicate attributes from different mergeSlices
				targetSlice = append(targetSlice, mergeAttribute)
			}
		}

	}
	return &targetSlice
}

func mergeDataSourceBoolAttribute(target, merge datasource.BoolAttribute) *datasource.BoolAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeDataSourceFloat64Attribute(target, merge datasource.Float64Attribute) *datasource.Float64Attribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeDataSourceInt64Attribute(target, merge datasource.Int64Attribute) *datasource.Int64Attribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeDataSourceNumberAttribute(target, merge datasource.NumberAttribute) *datasource.NumberAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeDataSourceStringAttribute(target, merge datasource.StringAttribute) *datasource.StringAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeDataSourceSingleNestedAttribute(target, merge datasource.SingleNestedAttribute) *datasource.SingleNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.Attributes = *MergeDataSourceAttributes(target.Attributes, merge.Attributes)

	return &target
}

func mergeDataSourceListNestedAttribute(target, merge datasource.ListNestedAttribute) *datasource.ListNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.NestedObject.Attributes = *MergeDataSourceAttributes(target.NestedObject.Attributes, merge.NestedObject.Attributes)

	return &target
}

func mergeDataSourceMapNestedAttribute(target, merge datasource.MapNestedAttribute) *datasource.MapNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.NestedObject.Attributes = *MergeDataSourceAttributes(target.NestedObject.Attributes, merge.NestedObject.Attributes)

	return &target
}

func mergeDataSourceSetNestedAttribute(target, merge datasource.SetNestedAttribute) *datasource.SetNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.NestedObject.Attributes = *MergeDataSourceAttributes(target.NestedObject.Attributes, merge.NestedObject.Attributes)

	return &target
}

func mergeDataSourceListAttribute(target, merge datasource.ListAttribute) *datasource.ListAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.ElementType = mergeElementType(target.ElementType, merge.ElementType)

	return &target
}

func mergeDataSourceMapAttribute(target, merge datasource.MapAttribute) *datasource.MapAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.ElementType = mergeElementType(target.ElementType, merge.ElementType)

	return &target
}

func mergeDataSourceSetAttribute(target, merge datasource.SetAttribute) *datasource.SetAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.ElementType = mergeElementType(target.ElementType, merge.ElementType)

	return &target
}
