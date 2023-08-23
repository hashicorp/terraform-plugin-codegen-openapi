// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type ResourceStringAttribute struct {
	resource.StringAttribute

	Name string
}

func (a *ResourceStringAttribute) GetName() string {
	return a.Name
}

func (a *ResourceStringAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	stringAttribute, ok := mergeAttribute.(*ResourceStringAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = stringAttribute.Description
	}

	return a
}

func (a *ResourceStringAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:   a.Name,
		String: &a.StringAttribute,
	}
}
