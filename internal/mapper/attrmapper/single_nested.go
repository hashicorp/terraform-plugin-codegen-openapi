// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type ResourceSingleNestedAttribute struct {
	resource.SingleNestedAttribute

	Name       string
	Attributes ResourceAttributes
}

func (a *ResourceSingleNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSingleNestedAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	singleNestedAttribute, ok := mergeAttribute.(*ResourceSingleNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = singleNestedAttribute.Description
	}
	a.Attributes = a.Attributes.Merge(singleNestedAttribute.Attributes)

	return a
}

func (a *ResourceSingleNestedAttribute) ToSpec() resource.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return resource.Attribute{
		Name:         a.Name,
		SingleNested: &a.SingleNestedAttribute,
	}
}
