// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
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

func (a *ResourceFloat64Attribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	float64Attribute, ok := mergeAttribute.(*ResourceFloat64Attribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = float64Attribute.Description
	}

	return a, nil
}

func (a *ResourceFloat64Attribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	cor, err := ApplyComputedOptionalRequiredOverride(override.ComputedOptionalRequired)
	if err != nil {
		return nil, err
	}
	a.ComputedOptionalRequired = cor

	return a, nil
}

func (a *ResourceFloat64Attribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:    util.TerraformIdentifier(a.Name),
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

func (a *DataSourceFloat64Attribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	float64Attribute, ok := mergeAttribute.(*DataSourceFloat64Attribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = float64Attribute.Description
	}

	return a, nil
}

func (a *DataSourceFloat64Attribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceFloat64Attribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:    util.TerraformIdentifier(a.Name),
		Float64: &a.Float64Attribute,
	}
}

type ProviderFloat64Attribute struct {
	provider.Float64Attribute

	Name string
}

func (a *ProviderFloat64Attribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name:    util.TerraformIdentifier(a.Name),
		Float64: &a.Float64Attribute,
	}
}
