package mapper

import (
	"github/hashicorp/terraform-providers-code-generator-openapi/internal/ir"
)

// mainSlice takes priority in the merge, will have each subsequent mergeAttributeSlice applied in sequence
// - No re-ordering of the mainSlice is done, so will append new attributes as they are encountered
func deepMergeAttributes(mainSlice []ir.ResourceAttribute, mergeAttributeSlices ...[]ir.ResourceAttribute) *[]ir.ResourceAttribute {
	for _, attributeSlice := range mergeAttributeSlices {

		for _, compareAttribute := range attributeSlice {
			isNewAttribute := true

			for mainIndex, mainAttribute := range mainSlice {
				if mainAttribute.Name == compareAttribute.Name {
					// Handle types that require nested merging
					if mainAttribute.SingleNested != nil && compareAttribute.SingleNested != nil {
						mergedAttributes := deepMergeAttributes(mainAttribute.SingleNested.Attributes, compareAttribute.SingleNested.Attributes)
						mainSlice[mainIndex].SingleNested.Attributes = *mergedAttributes
					} else if mainAttribute.ListNested != nil && compareAttribute.ListNested != nil {
						mergedAttributes := deepMergeAttributes(mainAttribute.ListNested.NestedObject.Attributes, compareAttribute.ListNested.NestedObject.Attributes)
						mainSlice[mainIndex].ListNested.NestedObject.Attributes = *mergedAttributes
					} else if mainAttribute.List != nil && compareAttribute.List != nil {
						mergedElementType := deepMergeElementType(&mainAttribute.List.ElementType, &compareAttribute.List.ElementType)
						mainSlice[mainIndex].List.ElementType = *mergedElementType
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

func deepMergeElementType(mainElementType *ir.ElementType, mergeElementType *ir.ElementType) *ir.ElementType {
	if mainElementType.List != nil && mergeElementType.List != nil {
		mainElementType.List.ElementType = deepMergeElementType(mainElementType.List.ElementType, mergeElementType.List.ElementType)
	} else if mainElementType.Object != nil && mergeElementType.Object != nil {
		objectElemTypes := deepMergeObjectElementTypes(mainElementType.Object, mergeElementType.Object)
		mainElementType.Object = objectElemTypes
	}

	return mainElementType
}

func deepMergeObjectElementTypes(mainObject []ir.ObjectElement, mergeObject []ir.ObjectElement) []ir.ObjectElement {
	for _, compareElemType := range mergeObject {
		isNewElemType := true

		for mainIndex, mainElemType := range mainObject {
			if mainElemType.Name == compareElemType.Name {
				mergedElementType := deepMergeElementType(mainElemType.ElementType, compareElemType.ElementType)
				mainObject[mainIndex].ElementType = mergedElementType

				isNewElemType = false
				break
			}
		}

		if isNewElemType {
			mainObject = append(mainObject, compareElemType)
		}
	}

	return mainObject
}
