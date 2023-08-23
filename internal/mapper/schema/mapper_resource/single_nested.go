// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperSingleNestedAttribute struct {
	resource.SingleNestedAttribute

	Name       string
	Attributes MapperAttributes
}

func (a *MapperSingleNestedAttribute) GetName() string {
	return a.Name
}

func (a *MapperSingleNestedAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	singleNestedAttribute, ok := mergeAttribute.(*MapperSingleNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = singleNestedAttribute.Description
	}
	a.Attributes = a.Attributes.Merge(singleNestedAttribute.Attributes)

	return a
}

func (a *MapperSingleNestedAttribute) ToSpec() resource.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return resource.Attribute{
		Name:         a.Name,
		SingleNested: &a.SingleNestedAttribute,
	}
}
