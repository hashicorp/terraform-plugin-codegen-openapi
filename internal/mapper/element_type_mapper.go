package mapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

func mergeElementType(mainElementType schema.ElementType, secondElementType schema.ElementType) schema.ElementType {
	if mainElementType.List != nil && secondElementType.List != nil {
		mainElementType.List.ElementType = mergeElementType(mainElementType.List.ElementType, secondElementType.List.ElementType)
	} else if mainElementType.Object != nil && secondElementType.Object != nil {
		objectElemTypes := mergeObjectAttributeTypes(mainElementType.Object, secondElementType.Object)
		mainElementType.Object = objectElemTypes
	}

	return mainElementType
}

func mergeObjectAttributeTypes(mainObject []schema.ObjectAttributeType, mergeObject []schema.ObjectAttributeType) []schema.ObjectAttributeType {
	for _, compareAttrType := range mergeObject {
		isNewElemType := true

		for mainIndex, mainAttrType := range mainObject {
			if mainAttrType.Name == compareAttrType.Name {
				// TODO: update here for new object attribute structure
				mergedElementType := mergeElementType(util.ConvertToElementType(mainAttrType), util.ConvertToElementType(compareAttrType))
				mainObject[mainIndex] = util.ConvertToAttributeType(mainAttrType.Name, mergedElementType)

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
