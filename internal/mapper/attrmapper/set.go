// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceSetAttribute struct {
	resource.SetAttribute

	Name string
}

func (a *ResourceSetAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSetAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	setAttribute, ok := mergeAttribute.(*ResourceSetAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, setAttribute.ElementType)

	return a
}

func (a *ResourceSetAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		Set:  &a.SetAttribute,
	}
}

type DataSourceSetAttribute struct {
	datasource.SetAttribute

	Name string
}

func (a *DataSourceSetAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceSetAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	setAttribute, ok := mergeAttribute.(*DataSourceSetAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, setAttribute.ElementType)

	return a
}

func (a *DataSourceSetAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: a.Name,
		Set:  &a.SetAttribute,
	}
}
