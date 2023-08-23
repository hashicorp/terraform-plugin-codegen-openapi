// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type ResourceFloat64Attribute struct {
	resource.Float64Attribute

	Name string
}

func (a *ResourceFloat64Attribute) GetName() string {
	return a.Name
}

func (a *ResourceFloat64Attribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	float64Attribute, ok := mergeAttribute.(*ResourceFloat64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = float64Attribute.Description
	}

	return a
}

func (a *ResourceFloat64Attribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:    a.Name,
		Float64: &a.Float64Attribute,
	}
}
