// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceListNestedAttribute struct {
	resource.ListNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceListNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceListNestedAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	listNestedAttribute, ok := mergeAttribute.(*ResourceListNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(listNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *ResourceListNestedAttribute) ToSpec() resource.Attribute {
	a.ListNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:       a.Name,
		ListNested: &a.ListNestedAttribute,
	}
}

type DataSourceListNestedAttribute struct {
	datasource.ListNestedAttribute

	Name         string
	NestedObject DataSourceNestedAttributeObject
}

func (a *DataSourceListNestedAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceListNestedAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	listNestedAttribute, ok := mergeAttribute.(*DataSourceListNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(listNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *DataSourceListNestedAttribute) ToSpec() datasource.Attribute {
	a.ListNestedAttribute.NestedObject = datasource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return datasource.Attribute{
		Name:       a.Name,
		ListNested: &a.ListNestedAttribute,
	}
}

type ProviderListNestedAttribute struct {
	provider.ListNestedAttribute

	Name         string
	NestedObject ProviderNestedAttributeObject
}

func (a *ProviderListNestedAttribute) ToSpec() provider.Attribute {
	a.ListNestedAttribute.NestedObject = provider.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return provider.Attribute{
		Name:       a.Name,
		ListNested: &a.ListNestedAttribute,
	}
}
