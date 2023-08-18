package merge

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

// MergeResourceAttributes takes a 'targetSlice' of resource attributes, which all other 'mergeSlices' of attributes
// are merged to. Attributes with the same name will be merged together and certain properties, like 'description', will be
// populated based on priority (targetSlice having the highest priority, then mergeSlices[0], mergeSlices[1], etc).
//
// All attributes not present in targetSlice will be appended to the end.
func MergeResourceAttributes(targetSlice []resource.Attribute, mergeSlices ...[]resource.Attribute) *[]resource.Attribute {
	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.Name == mergeAttribute.Name {
					// Handle primitive attribute type
					// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#primitive-attribute-types
					if targetAttribute.Bool != nil && mergeAttribute.Bool != nil {
						targetSlice[i].Bool = mergeResourceBoolAttribute(*targetAttribute.Bool, *mergeAttribute.Bool)
					} else if targetAttribute.Float64 != nil && mergeAttribute.Float64 != nil {
						targetSlice[i].Float64 = mergeResourceFloat64Attribute(*targetAttribute.Float64, *mergeAttribute.Float64)
					} else if targetAttribute.Int64 != nil && mergeAttribute.Int64 != nil {
						targetSlice[i].Int64 = mergeResourceInt64Attribute(*targetAttribute.Int64, *mergeAttribute.Int64)
					} else if targetAttribute.Number != nil && mergeAttribute.Number != nil {
						targetSlice[i].Number = mergeResourceNumberAttribute(*targetAttribute.Number, *mergeAttribute.Number)
					} else if targetAttribute.String != nil && mergeAttribute.String != nil {
						targetSlice[i].String = mergeResourceStringAttribute(*targetAttribute.String, *mergeAttribute.String)
					}

					// Handle nested attribute type
					// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#nested-attribute-types
					if targetAttribute.SingleNested != nil && mergeAttribute.SingleNested != nil {
						targetSlice[i].SingleNested = mergeResourceSingleNestedAttribute(*targetAttribute.SingleNested, *mergeAttribute.SingleNested)
					} else if targetAttribute.ListNested != nil && mergeAttribute.ListNested != nil {
						targetSlice[i].ListNested = mergeResourceListNestedAttribute(*targetAttribute.ListNested, *mergeAttribute.ListNested)
					} else if targetAttribute.MapNested != nil && mergeAttribute.MapNested != nil {
						targetSlice[i].MapNested = mergeResourceMapNestedAttribute(*targetAttribute.MapNested, *mergeAttribute.MapNested)
					} else if targetAttribute.SetNested != nil && mergeAttribute.SetNested != nil {
						targetSlice[i].SetNested = mergeResourceSetNestedAttribute(*targetAttribute.SetNested, *mergeAttribute.SetNested)
					}

					// Handle collection attribute type
					// https://developer.hashicorp.com/terraform/plugin/framework/handling-data/attributes#collection-attribute-types
					if targetAttribute.List != nil && mergeAttribute.List != nil {
						targetSlice[i].List = mergeResourceListAttribute(*targetAttribute.List, *mergeAttribute.List)
					} else if targetAttribute.Map != nil && mergeAttribute.Map != nil {
						targetSlice[i].Map = mergeResourceMapAttribute(*targetAttribute.Map, *mergeAttribute.Map)
					} else if targetAttribute.Set != nil && mergeAttribute.Set != nil {
						targetSlice[i].Set = mergeResourceSetAttribute(*targetAttribute.Set, *mergeAttribute.Set)
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

func mergeResourceBoolAttribute(target, merge resource.BoolAttribute) *resource.BoolAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeResourceFloat64Attribute(target, merge resource.Float64Attribute) *resource.Float64Attribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeResourceInt64Attribute(target, merge resource.Int64Attribute) *resource.Int64Attribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeResourceNumberAttribute(target, merge resource.NumberAttribute) *resource.NumberAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeResourceStringAttribute(target, merge resource.StringAttribute) *resource.StringAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	return &target
}

func mergeResourceSingleNestedAttribute(target, merge resource.SingleNestedAttribute) *resource.SingleNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.Attributes = *MergeResourceAttributes(target.Attributes, merge.Attributes)

	return &target
}

func mergeResourceListNestedAttribute(target, merge resource.ListNestedAttribute) *resource.ListNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.NestedObject.Attributes = *MergeResourceAttributes(target.NestedObject.Attributes, merge.NestedObject.Attributes)

	return &target
}

func mergeResourceMapNestedAttribute(target, merge resource.MapNestedAttribute) *resource.MapNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.NestedObject.Attributes = *MergeResourceAttributes(target.NestedObject.Attributes, merge.NestedObject.Attributes)

	return &target
}

func mergeResourceSetNestedAttribute(target, merge resource.SetNestedAttribute) *resource.SetNestedAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.NestedObject.Attributes = *MergeResourceAttributes(target.NestedObject.Attributes, merge.NestedObject.Attributes)

	return &target
}

func mergeResourceListAttribute(target, merge resource.ListAttribute) *resource.ListAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.ElementType = mergeElementType(target.ElementType, merge.ElementType)

	return &target
}

func mergeResourceMapAttribute(target, merge resource.MapAttribute) *resource.MapAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.ElementType = mergeElementType(target.ElementType, merge.ElementType)

	return &target
}

func mergeResourceSetAttribute(target, merge resource.SetAttribute) *resource.SetAttribute {
	if target.Description == nil || *target.Description == "" {
		target.Description = merge.Description
	}

	target.ElementType = mergeElementType(target.ElementType, merge.ElementType)

	return &target
}
