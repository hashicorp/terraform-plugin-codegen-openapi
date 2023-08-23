// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type ResourceSetNestedAttribute struct {
	resource.SetNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceSetNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSetNestedAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	setNestedAttribute, ok := mergeAttribute.(*ResourceSetNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(setNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *ResourceSetNestedAttribute) ToSpec() resource.Attribute {
	a.SetNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      a.Name,
		SetNested: &a.SetNestedAttribute,
	}
}
