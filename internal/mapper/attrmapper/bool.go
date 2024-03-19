// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
)

type ResourceBoolAttribute struct {
	resource.BoolAttribute

	Name string
}

func (a *ResourceBoolAttribute) GetName() string {
	return a.Name
}

func (a *ResourceBoolAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	boolAttribute, ok := mergeAttribute.(*ResourceBoolAttribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = boolAttribute.Description
	}

	return a, nil
}

func (a *ResourceBoolAttribute) ApplyOverride(override config.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *ResourceBoolAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Bool: &a.BoolAttribute,
	}
}

type DataSourceBoolAttribute struct {
	datasource.BoolAttribute

	Name string
}

func (a *DataSourceBoolAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceBoolAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	boolAttribute, ok := mergeAttribute.(*DataSourceBoolAttribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = boolAttribute.Description
	}

	return a, nil
}

func (a *DataSourceBoolAttribute) ApplyOverride(override config.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceBoolAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Bool: &a.BoolAttribute,
	}
}

type ProviderBoolAttribute struct {
	provider.BoolAttribute

	Name string
}

func (a *ProviderBoolAttribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Bool: &a.BoolAttribute,
	}
}
