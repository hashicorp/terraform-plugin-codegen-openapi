// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceSetNestedAttribute struct {
	resource.SetNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceSetNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSetNestedAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	setNestedAttribute, ok := mergeAttribute.(*ResourceSetNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(setNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *ResourceSetNestedAttribute) ToSpec() resource.Attribute {
	a.SetNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      a.Name,
		SetNested: &a.SetNestedAttribute,
	}
}

type DataSourceSetNestedAttribute struct {
	datasource.SetNestedAttribute

	Name         string
	NestedObject DataSourceNestedAttributeObject
}

func (a *DataSourceSetNestedAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceSetNestedAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	setNestedAttribute, ok := mergeAttribute.(*DataSourceSetNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(setNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *DataSourceSetNestedAttribute) ToSpec() datasource.Attribute {
	a.SetNestedAttribute.NestedObject = datasource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return datasource.Attribute{
		Name:      a.Name,
		SetNested: &a.SetNestedAttribute,
	}
}

type ProviderSetNestedAttribute struct {
	provider.SetNestedAttribute

	Name         string
	NestedObject ProviderNestedAttributeObject
}

func (a *ProviderSetNestedAttribute) ToSpec() provider.Attribute {
	a.SetNestedAttribute.NestedObject = provider.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return provider.Attribute{
		Name:      a.Name,
		SetNested: &a.SetNestedAttribute,
	}
}