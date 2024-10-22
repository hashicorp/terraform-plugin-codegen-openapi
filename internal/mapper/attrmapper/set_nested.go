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

type ResourceSetNestedAttribute struct {
	resource.SetNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceSetNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSetNestedAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	setNestedAttribute, ok := mergeAttribute.(*ResourceSetNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setNestedAttribute.Description
	}
	a.NestedObject.Attributes, _ = a.NestedObject.Attributes.Merge(setNestedAttribute.NestedObject.Attributes)

	return a, nil
}

func (a *ResourceSetNestedAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *ResourceSetNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (ResourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *ResourceSetNestedAttribute) NestedMerge(path []string, attribute ResourceAttribute, intermediateComputability schema.ComputedOptionalRequired) (ResourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.MergeAttribute(path, attribute, intermediateComputability)
	if err == nil {
		a.ComputedOptionalRequired = intermediateComputability
	}

	return a, err
}

func (a *ResourceSetNestedAttribute) ToSpec() resource.Attribute {
	a.SetNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:      util.TerraformIdentifier(a.Name),
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

func (a *DataSourceSetNestedAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	setNestedAttribute, ok := mergeAttribute.(*DataSourceSetNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setNestedAttribute.Description
	}
	a.NestedObject.Attributes, _ = a.NestedObject.Attributes.Merge(setNestedAttribute.NestedObject.Attributes)

	return a, nil
}

func (a *DataSourceSetNestedAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceSetNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (DataSourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *DataSourceSetNestedAttribute) NestedMerge(path []string, attribute DataSourceAttribute, intermediateComputability schema.ComputedOptionalRequired) (DataSourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.MergeAttribute(path, attribute, intermediateComputability)
	if err == nil {
		a.ComputedOptionalRequired = intermediateComputability
	}

	return a, err
}

func (a *DataSourceSetNestedAttribute) ToSpec() datasource.Attribute {
	a.SetNestedAttribute.NestedObject = datasource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return datasource.Attribute{
		Name:      util.TerraformIdentifier(a.Name),
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
		Name:      util.TerraformIdentifier(a.Name),
		SetNested: &a.SetNestedAttribute,
	}
}
