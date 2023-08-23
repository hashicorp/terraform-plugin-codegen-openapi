// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceNumberAttribute struct {
	resource.NumberAttribute

	Name string
}

func (a *ResourceNumberAttribute) GetName() string {
	return a.Name
}

func (a *ResourceNumberAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	numberAttribute, ok := mergeAttribute.(*ResourceNumberAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = numberAttribute.Description
	}

	return a
}

func (a *ResourceNumberAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:   a.Name,
		Number: &a.NumberAttribute,
	}
}

type DataSourceNumberAttribute struct {
	datasource.NumberAttribute

	Name string
}

func (a *DataSourceNumberAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceNumberAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	numberAttribute, ok := mergeAttribute.(*DataSourceNumberAttribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = numberAttribute.Description
	}

	return a
}

func (a *DataSourceNumberAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:   a.Name,
		Number: &a.NumberAttribute,
	}
}
