// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceAttribute interface {
	GetName() string
	Merge(ResourceAttribute) (ResourceAttribute, error)
	ToSpec() resource.Attribute
}

type ResourceAttributes []ResourceAttribute

func (targetSlice ResourceAttributes) Merge(mergeSlices ...ResourceAttributes) (ResourceAttributes, error) {
	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.GetName() == mergeAttribute.GetName() {
					// TODO: determine how to surface this error
					targetSlice[i], _ = targetAttribute.Merge(mergeAttribute)

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

	return targetSlice, nil
}

func (attributes ResourceAttributes) ToSpec() []resource.Attribute {
	specAttributes := make([]resource.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}

type DataSourceAttribute interface {
	GetName() string
	Merge(DataSourceAttribute) (DataSourceAttribute, error)
	ToSpec() datasource.Attribute
}

type DataSourceAttributes []DataSourceAttribute

func (targetSlice DataSourceAttributes) Merge(mergeSlices ...DataSourceAttributes) (DataSourceAttributes, error) {
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

func (attributes DataSourceAttributes) ToSpec() []datasource.Attribute {
	specAttributes := make([]datasource.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}

type ProviderAttribute interface {
	ToSpec() provider.Attribute
}

type ProviderAttributes []ProviderAttribute

func (attributes ProviderAttributes) ToSpec() []provider.Attribute {
	specAttributes := make([]provider.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}
