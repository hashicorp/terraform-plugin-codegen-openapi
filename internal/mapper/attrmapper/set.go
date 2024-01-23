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

type ResourceSetAttribute struct {
	resource.SetAttribute

	Name string
}

func (a *ResourceSetAttribute) GetName() string {
	return a.Name
}

func (a *ResourceSetAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	setAttribute, ok := mergeAttribute.(*ResourceSetAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, setAttribute.ElementType)

	return a, nil
}

func (a *ResourceSetAttribute) ApplyOverride(override config.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *ResourceSetAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Set:  &a.SetAttribute,
	}
}

type DataSourceSetAttribute struct {
	datasource.SetAttribute

	Name string
}

func (a *DataSourceSetAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceSetAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	setAttribute, ok := mergeAttribute.(*DataSourceSetAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = setAttribute.Description
	}
	a.ElementType = mergeElementType(a.ElementType, setAttribute.ElementType)

	return a, nil
}

func (a *DataSourceSetAttribute) ApplyOverride(override config.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceSetAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Set:  &a.SetAttribute,
	}
}

type ProviderSetAttribute struct {
	provider.SetAttribute

	Name string
}

func (a *ProviderSetAttribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name: util.TerraformIdentifier(a.Name),
		Set:  &a.SetAttribute,
	}
}
