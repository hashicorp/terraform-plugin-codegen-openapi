// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceFloat64Attribute struct {
	resource.Float64Attribute

	Name string
}

func (a *ResourceFloat64Attribute) GetName() string {
	return a.Name
}

func (a *ResourceFloat64Attribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	float64Attribute, ok := mergeAttribute.(*ResourceFloat64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = float64Attribute.Description
	}

	return a
}

func (a *ResourceFloat64Attribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:    a.Name,
		Float64: &a.Float64Attribute,
	}
}

type DataSourceFloat64Attribute struct {
	datasource.Float64Attribute

	Name string
}

func (a *DataSourceFloat64Attribute) GetName() string {
	return a.Name
}

func (a *DataSourceFloat64Attribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	float64Attribute, ok := mergeAttribute.(*DataSourceFloat64Attribute)
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = float64Attribute.Description
	}

	return a
}

func (a *DataSourceFloat64Attribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:    a.Name,
		Float64: &a.Float64Attribute,
	}
}

type ProviderFloat64Attribute struct {
	provider.Float64Attribute

	Name string
}

func (a *ProviderFloat64Attribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name:    a.Name,
		Float64: &a.Float64Attribute,
	}
}
