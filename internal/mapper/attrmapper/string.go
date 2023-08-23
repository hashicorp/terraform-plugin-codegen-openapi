// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

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

type DataSourceStringAttribute struct {
	datasource.StringAttribute

	Name string
}

func (a *DataSourceStringAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceStringAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	stringAttribute, ok := mergeAttribute.(*DataSourceStringAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = stringAttribute.Description
	}

	return a
}

func (a *DataSourceStringAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:   a.Name,
		String: &a.StringAttribute,
	}
}
