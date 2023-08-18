package merge

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

// mainSlice takes priority in the merge, will have each subsequent mergeAttributeSlice applied in sequence
// - No re-ordering of the mainSlice is done, so will append new attributes as they are encountered
func MergeResourceAttributes(mainSlice []resource.Attribute, mergeAttributeSlices ...[]resource.Attribute) *[]resource.Attribute {
	// loop through each slice of attributes by priority (order)
	for _, attributeSlice := range mergeAttributeSlices {

		// loop through each attribute
		for _, compareAttribute := range attributeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			// Loop through the main slice's attributes and try to find a match
			for mainIndex, mainAttribute := range mainSlice {
				// If the attribute names match, we need to dive deeper
				if mainAttribute.Name == compareAttribute.Name {
					// Handle types that require nested merging
					if mainAttribute.SingleNested != nil && compareAttribute.SingleNested != nil {
						mergedAttributes := MergeResourceAttributes(mainAttribute.SingleNested.Attributes, compareAttribute.SingleNested.Attributes)
						mainSlice[mainIndex].SingleNested.Attributes = *mergedAttributes
					} else if mainAttribute.ListNested != nil && compareAttribute.ListNested != nil {
						mergedAttributes := MergeResourceAttributes(mainAttribute.ListNested.NestedObject.Attributes, compareAttribute.ListNested.NestedObject.Attributes)
						mainSlice[mainIndex].ListNested.NestedObject.Attributes = *mergedAttributes
					} else if mainAttribute.List != nil && compareAttribute.List != nil {
						mergedElementType := mergeElementType(mainAttribute.List.ElementType, compareAttribute.List.ElementType)
						mainSlice[mainIndex].List.ElementType = mergedElementType
					}

					isNewAttribute = false
					break
				}
			}

			if isNewAttribute {
				// Add this back to the original slice to avoid adding duplicate attributes from different mergeAttributeSlices
				mainSlice = append(mainSlice, compareAttribute)
			}
		}

	}
	return &mainSlice
}

// mainSlice takes priority in the merge, will have each subsequent mergeAttributeSlice applied in sequence
// - No re-ordering of the mainSlice is done, so will append new attributes as they are encountered
func MergeDataSourceAttributes(mainSlice []datasource.Attribute, mergeAttributeSlices ...[]datasource.Attribute) *[]datasource.Attribute {
	// loop through each slice of attributes by priority (order)
	for _, attributeSlice := range mergeAttributeSlices {

		// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
		for _, compareAttribute := range attributeSlice {
			isNewAttribute := true

			// Loop through the main slice's attributes and try to find a match
			for mainIndex, mainAttribute := range mainSlice {
				// If the attribute names match, we need to dive deeper
				if mainAttribute.Name == compareAttribute.Name {
					// Handle types that require nested merging
					if mainAttribute.SingleNested != nil && compareAttribute.SingleNested != nil {
						mergedAttributes := MergeDataSourceAttributes(mainAttribute.SingleNested.Attributes, compareAttribute.SingleNested.Attributes)
						mainSlice[mainIndex].SingleNested.Attributes = *mergedAttributes
					} else if mainAttribute.ListNested != nil && compareAttribute.ListNested != nil {
						mergedAttributes := MergeDataSourceAttributes(mainAttribute.ListNested.NestedObject.Attributes, compareAttribute.ListNested.NestedObject.Attributes)
						mainSlice[mainIndex].ListNested.NestedObject.Attributes = *mergedAttributes
					} else if mainAttribute.List != nil && compareAttribute.List != nil {
						mergedElementType := mergeElementType(mainAttribute.List.ElementType, compareAttribute.List.ElementType)
						mainSlice[mainIndex].List.ElementType = mergedElementType
					}

					isNewAttribute = false
					break
				}
			}

			if isNewAttribute {
				// Add this back to the original slice to avoid adding duplicate attributes from different mergeAttributeSlices
				mainSlice = append(mainSlice, compareAttribute)
			}
		}

	}
	return &mainSlice
}

func mergeElementType(mainElementType schema.ElementType, secondElementType schema.ElementType) schema.ElementType {
	if mainElementType.List != nil && secondElementType.List != nil {
		mainElementType.List.ElementType = mergeElementType(mainElementType.List.ElementType, secondElementType.List.ElementType)
	} else if mainElementType.Object != nil && secondElementType.Object != nil {
		objectElemTypes := mergeObjectAttributeTypes(mainElementType.Object.AttributeTypes, secondElementType.Object.AttributeTypes)
		mainElementType.Object.AttributeTypes = objectElemTypes
	}

	return mainElementType
}

func mergeObjectAttributeTypes(mainObject []schema.ObjectAttributeType, mergeObject []schema.ObjectAttributeType) []schema.ObjectAttributeType {
	for _, compareAttrType := range mergeObject {
		isNewElemType := true

		for mainIndex, mainAttrType := range mainObject {
			if mainAttrType.Name == compareAttrType.Name {
				mergedElementType := mergeElementType(util.CreateElementType(mainAttrType), util.CreateElementType(compareAttrType))
				mainObject[mainIndex] = util.CreateObjectAttributeType(mainAttrType.Name, mergedElementType)

				isNewElemType = false
				break
			}
		}

		if isNewElemType {
			mainObject = append(mainObject, compareAttrType)
		}
	}

	return mainObject
}
