// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceMapNestedAttribute struct {
	resource.MapNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceMapNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceMapNestedAttribute) Merge(mergeAttribute ResourceAttribute) ResourceAttribute {
	mapNestedAttribute, ok := mergeAttribute.(*ResourceMapNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(mapNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *ResourceMapNestedAttribute) ToSpec() resource.Attribute {
	a.MapNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      a.Name,
		MapNested: &a.MapNestedAttribute,
	}
}

type DataSourceMapNestedAttribute struct {
	datasource.MapNestedAttribute

	Name         string
	NestedObject DataSourceNestedAttributeObject
}

func (a *DataSourceMapNestedAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceMapNestedAttribute) Merge(mergeAttribute DataSourceAttribute) DataSourceAttribute {
	mapNestedAttribute, ok := mergeAttribute.(*DataSourceMapNestedAttribute)
	if !ok {
		return a
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapNestedAttribute.Description
	}
	a.NestedObject.Attributes = a.NestedObject.Attributes.Merge(mapNestedAttribute.NestedObject.Attributes)

	return a
}

func (a *DataSourceMapNestedAttribute) ToSpec() datasource.Attribute {
	a.MapNestedAttribute.NestedObject = datasource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return datasource.Attribute{
		Name:      a.Name,
		MapNested: &a.MapNestedAttribute,
	}
}

type ProviderMapNestedAttribute struct {
	provider.MapNestedAttribute

	Name         string
	NestedObject ProviderNestedAttributeObject
}

func (a *ProviderMapNestedAttribute) ToSpec() provider.Attribute {
	a.MapNestedAttribute.NestedObject = provider.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return provider.Attribute{
		Name:      a.Name,
		MapNested: &a.MapNestedAttribute,
	}
}
