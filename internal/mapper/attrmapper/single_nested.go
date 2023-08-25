// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceSingleNestedAttribute struct {
	resource.SingleNestedAttribute

	Name       string
	Attributes ResourceAttributes
}

func (a *ResourceSingleNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSingleNestedAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	singleNestedAttribute, ok := mergeAttribute.(*ResourceSingleNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = singleNestedAttribute.Description
	}
	a.Attributes, _ = a.Attributes.Merge(singleNestedAttribute.Attributes)

	return a, nil
}

func (a *ResourceSingleNestedAttribute) ToSpec() resource.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return resource.Attribute{
		Name:         a.Name,
		SingleNested: &a.SingleNestedAttribute,
	}
}

type DataSourceSingleNestedAttribute struct {
	datasource.SingleNestedAttribute

	Name       string
	Attributes DataSourceAttributes
}

func (a *DataSourceSingleNestedAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceSingleNestedAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	singleNestedAttribute, ok := mergeAttribute.(*DataSourceSingleNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = singleNestedAttribute.Description
	}
	a.Attributes, _ = a.Attributes.Merge(singleNestedAttribute.Attributes)

	return a, nil
}

func (a *DataSourceSingleNestedAttribute) ToSpec() datasource.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return datasource.Attribute{
		Name:         a.Name,
		SingleNested: &a.SingleNestedAttribute,
	}
}

type ProviderSingleNestedAttribute struct {
	provider.SingleNestedAttribute

	Name       string
	Attributes ProviderAttributes
}

func (a *ProviderSingleNestedAttribute) ToSpec() provider.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return provider.Attribute{
		Name:         a.Name,
		SingleNested: &a.SingleNestedAttribute,
	}
}
