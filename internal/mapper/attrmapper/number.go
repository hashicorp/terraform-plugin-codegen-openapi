// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceNumberAttribute struct {
	resource.NumberAttribute

	Name string
}

func (a *ResourceNumberAttribute) GetName() string {
	return a.Name
}

func (a *ResourceNumberAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	numberAttribute, ok := mergeAttribute.(*ResourceNumberAttribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = numberAttribute.Description
	}

	return a, nil
}

func (a *ResourceNumberAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
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

func (a *DataSourceNumberAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	numberAttribute, ok := mergeAttribute.(*DataSourceNumberAttribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = numberAttribute.Description
	}

	return a, nil
}

func (a *DataSourceNumberAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceNumberAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:   a.Name,
		Number: &a.NumberAttribute,
	}
}

type ProviderNumberAttribute struct {
	provider.NumberAttribute

	Name string
}

func (a *ProviderNumberAttribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name:   a.Name,
		Number: &a.NumberAttribute,
	}
}
