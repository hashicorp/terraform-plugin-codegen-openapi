// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceListAttribute struct {
	resource.ListAttribute

	Name string
}

func (a *ResourceListAttribute) GetName() string {
	return a.Name
}

func (a *ResourceListAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	listAttribute, ok := mergeAttribute.(*ResourceListAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, listAttribute.ElementType)

	return a
}

func (a *ResourceListAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		List: &a.ListAttribute,
	}
}

type DataSourceListAttribute struct {
	datasource.ListAttribute

	Name string
}

func (a *DataSourceListAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceListAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	listAttribute, ok := mergeAttribute.(*DataSourceListAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, listAttribute.ElementType)

	return a
}

func (a *DataSourceListAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: a.Name,
		List: &a.ListAttribute,
	}
}
