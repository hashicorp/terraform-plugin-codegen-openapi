// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type ResourceNumberAttribute struct {
	resource.NumberAttribute

	Name string
}

func (a *ResourceNumberAttribute) GetName() string {
	return a.Name
}

func (a *ResourceNumberAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	numberAttribute, ok := mergeAttribute.(*ResourceNumberAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = numberAttribute.Description
	}

	return a
}

func (a *ResourceNumberAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:   a.Name,
		Number: &a.NumberAttribute,
	}
}
