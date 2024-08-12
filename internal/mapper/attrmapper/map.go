// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/raphaelfff/terraform-plugin-codegen-openapi/internal/mapper/util"
)

type ResourceMapAttribute struct {
	resource.MapAttribute

	Name string
}

func (a *ResourceMapAttribute) GetName() string {
	return a.Name
}

func (a *ResourceMapAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	mapAttribute, ok := mergeAttribute.(*ResourceMapAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, mapAttribute.ElementType)

	return a, nil
}

func (a *ResourceMapAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *ResourceMapAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Map:  &a.MapAttribute,
	}
}

type DataSourceMapAttribute struct {
	datasource.MapAttribute

	Name string
}

func (a *DataSourceMapAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceMapAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	mapAttribute, ok := mergeAttribute.(*DataSourceMapAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = mapAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, mapAttribute.ElementType)

	return a, nil
}

func (a *DataSourceMapAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceMapAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Map:  &a.MapAttribute,
	}
}

type ProviderMapAttribute struct {
	provider.MapAttribute

	Name string
}

func (a *ProviderMapAttribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Map:  &a.MapAttribute,
	}
}
