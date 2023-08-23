// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperInt64Attribute struct {
	resource.Int64Attribute

	Name string
}

func (a *MapperInt64Attribute) GetName() string {
	return a.Name
}

func (a *MapperInt64Attribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	int64Attribute, ok := mergeAttribute.(*MapperInt64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = int64Attribute.Description
	}

	return a
}

func (a *MapperInt64Attribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:  a.Name,
		Int64: &a.Int64Attribute,
	}
}
