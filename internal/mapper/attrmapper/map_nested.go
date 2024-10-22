// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

type ResourceMapNestedAttribute struct {
	resource.MapNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceMapNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceMapNestedAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	mapNestedAttribute, ok := mergeAttribute.(*ResourceMapNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapNestedAttribute.Description
	}
	a.NestedObject.Attributes, _ = a.NestedObject.Attributes.Merge(mapNestedAttribute.NestedObject.Attributes)

	return a, nil
}

func (a *ResourceMapNestedAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *ResourceMapNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (ResourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *ResourceMapNestedAttribute) NestedMerge(path []string, attribute ResourceAttribute, intermediateComputability schema.ComputedOptionalRequired) (ResourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.MergeAttribute(path, attribute, intermediateComputability)
	if err == nil {
		a.ComputedOptionalRequired = intermediateComputability
	}

	return a, err
}

func (a *ResourceMapNestedAttribute) ToSpec() resource.Attribute {
	a.MapNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      util.TerraformIdentifier(a.Name),
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

func (a *DataSourceMapNestedAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	mapNestedAttribute, ok := mergeAttribute.(*DataSourceMapNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapNestedAttribute.Description
	}
	a.NestedObject.Attributes, _ = a.NestedObject.Attributes.Merge(mapNestedAttribute.NestedObject.Attributes)

	return a, nil
}

func (a *DataSourceMapNestedAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceMapNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (DataSourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *DataSourceMapNestedAttribute) NestedMerge(path []string, attribute DataSourceAttribute, intermediateComputability schema.ComputedOptionalRequired) (DataSourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.MergeAttribute(path, attribute, intermediateComputability)
	if err == nil {
		a.ComputedOptionalRequired = intermediateComputability
	}

	return a, err
}

func (a *DataSourceMapNestedAttribute) ToSpec() datasource.Attribute {
	a.MapNestedAttribute.NestedObject = datasource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return datasource.Attribute{
		Name:      util.TerraformIdentifier(a.Name),
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
		Name:      util.TerraformIdentifier(a.Name),
		MapNested: &a.MapNestedAttribute,
	}
}
