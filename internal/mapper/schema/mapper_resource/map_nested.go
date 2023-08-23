// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mapper_resource

import "github.com/hashicorp/terraform-plugin-codegen-spec/resource"

type MapperMapNestedAttribute struct {
	resource.MapNestedAttribute

	Name         string
	NestedObject MapperNestedAttributeObject
}

func (a *MapperMapNestedAttribute) GetName() string {
	return a.Name
}

func (a *MapperMapNestedAttribute) Merge(mergeAttribute MapperAttribute) MapperAttribute {
	mapNestedAttribute, ok := mergeAttribute.(*MapperMapNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(mapNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *MapperMapNestedAttribute) ToSpec() resource.Attribute {
	a.MapNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      a.Name,
		MapNested: &a.MapNestedAttribute,
	}
}
