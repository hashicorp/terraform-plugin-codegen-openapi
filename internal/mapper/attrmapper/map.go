// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceMapAttribute struct {
	resource.MapAttribute

	Name string
}

func (a *ResourceMapAttribute) GetName() string {
	return a.Name
}

func (a *ResourceMapAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	mapAttribute, ok := mergeAttribute.(*ResourceMapAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, mapAttribute.ElementType)

	return a
}

func (a *ResourceMapAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		Map:  &a.MapAttribute,
	}
}

type DataSourceMapAttribute struct {
	datasource.MapAttribute

	Name string
}

func (a *DataSourceMapAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceMapAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	mapAttribute, ok := mergeAttribute.(*DataSourceMapAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, mapAttribute.ElementType)

	return a
}

func (a *DataSourceMapAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: a.Name,
		Map:  &a.MapAttribute,
	}
}
