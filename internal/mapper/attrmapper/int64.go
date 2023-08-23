// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceInt64Attribute struct {
	resource.Int64Attribute

	Name string
}

func (a *ResourceInt64Attribute) GetName() string {
	return a.Name
}

func (a *ResourceInt64Attribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	int64Attribute, ok := mergeAttribute.(*ResourceInt64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = int64Attribute.Description
	}

	return a
}

func (a *ResourceInt64Attribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:  a.Name,
		Int64: &a.Int64Attribute,
	}
}

type DataSourceInt64Attribute struct {
	datasource.Int64Attribute

	Name string
}

func (a *DataSourceInt64Attribute) GetName() string {
	return a.Name
}

func (a *DataSourceInt64Attribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	int64Attribute, ok := mergeAttribute.(*DataSourceInt64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = int64Attribute.Description
	}

	return a
}

func (a *DataSourceInt64Attribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:  a.Name,
		Int64: &a.Int64Attribute,
	}
}

type ProviderInt64Attribute struct {
	provider.Int64Attribute

	Name string
}

func (a *ProviderInt64Attribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name:  a.Name,
		Int64: &a.Int64Attribute,
	}
}
