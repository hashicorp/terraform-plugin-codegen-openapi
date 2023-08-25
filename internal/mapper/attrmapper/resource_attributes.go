// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceAttribute interface {
	GetName() string
	Merge(ResourceAttribute) (ResourceAttribute, error)
	ToSpec() resource.Attribute
}

type ResourceAttributes []ResourceAttribute

func (targetSlice ResourceAttributes) Merge(mergeSlices ...ResourceAttributes) (ResourceAttributes, error) {
	var errResult error

	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.GetName() == mergeAttribute.GetName() {
					mergedAttribute, err := targetAttribute.Merge(mergeAttribute)
					if err != nil {
						// TODO: consider how best to surface this error
						// Currently, if the merge fails we should just keep the original target attribute for now
						errResult = errors.Join(errResult, err)
					} else {
						targetSlice[i] = mergedAttribute
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

	return targetSlice, errResult
}

func (attributes ResourceAttributes) ToSpec() []resource.Attribute {
	specAttributes := make([]resource.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}
