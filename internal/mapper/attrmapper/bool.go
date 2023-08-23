// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceBoolAttribute struct {
	resource.BoolAttribute

	Name string
}

func (a *ResourceBoolAttribute) GetName() string {
	return a.Name
}

func (a *ResourceBoolAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	boolAttribute, ok := mergeAttribute.(*ResourceBoolAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = boolAttribute.Description
	}

	return a
}

func (a *ResourceBoolAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: a.Name,
		Bool: &a.BoolAttribute,
	}
}

type DataSourceBoolAttribute struct {
	datasource.BoolAttribute

	Name string
}

func (a *DataSourceBoolAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceBoolAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	boolAttribute, ok := mergeAttribute.(*DataSourceBoolAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = boolAttribute.Description
	}

	return a
}

func (a *DataSourceBoolAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: a.Name,
		Bool: &a.BoolAttribute,
	}
}
