// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceAttribute interface {
	GetName() string
	Merge(ResourceAttribute) ResourceAttribute
	ToSpec() resource.Attribute
}

type ResourceAttributes []ResourceAttribute

func (targetSlice ResourceAttributes) Merge(mergeSlices ...ResourceAttributes) ResourceAttributes {
	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.GetName() == mergeAttribute.GetName() {
					targetSlice[i] = targetAttribute.Merge(mergeAttribute)

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

	return targetSlice
}

func (attributes ResourceAttributes) ToSpec() []resource.Attribute {
	specAttributes := []resource.Attribute{}
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}

type DataSourceAttribute interface {
	GetName() string
	Merge(DataSourceAttribute) DataSourceAttribute
	ToSpec() datasource.Attribute
}

type DataSourceAttributes []DataSourceAttribute

func (targetSlice DataSourceAttributes) Merge(mergeSlices ...DataSourceAttributes) DataSourceAttributes {
	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.GetName() == mergeAttribute.GetName() {
					targetSlice[i] = targetAttribute.Merge(mergeAttribute)

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

	return targetSlice
}

func (attributes DataSourceAttributes) ToSpec() []datasource.Attribute {
	specAttributes := []datasource.Attribute{}
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
	specAttributes := []provider.Attribute{}
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}
