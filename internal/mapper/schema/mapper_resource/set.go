// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperSetAttribute struct {
	resource.SetAttribute

	Name string
}

func (a *MapperSetAttribute) GetName() string {
	return a.Name
}

func (a *MapperSetAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	setAttribute, ok := mergeAttribute.(*MapperSetAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, setAttribute.ElementType)

	return a
}

func (a *MapperSetAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		Set:  &a.SetAttribute,
	}
}
