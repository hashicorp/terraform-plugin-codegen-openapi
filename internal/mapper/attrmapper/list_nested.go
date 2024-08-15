// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/util"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

type ResourceListNestedAttribute struct {
	resource.ListNestedAttribute

	Name         string
	NestedObject ResourceNestedAttributeObject
}

func (a *ResourceListNestedAttribute) GetName() string {
	return a.Name
}

func (a *ResourceListNestedAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	listNestedAttribute, ok := mergeAttribute.(*ResourceListNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listNestedAttribute.Description
	}
	a.NestedObject.Attributes, _ = a.NestedObject.Attributes.Merge(listNestedAttribute.NestedObject.Attributes)

	return a, nil
}

func (a *ResourceListNestedAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
	a.Description = &override.Description

	switch override.ComputedOptionalRequired {
	case "": // No override
	case "computed":
		a.ComputedOptionalRequired = schema.Computed
	case "optional":
		a.ComputedOptionalRequired = schema.Optional
	case "required":
		a.ComputedOptionalRequired = schema.Required
	case "computed_optional":
		a.ComputedOptionalRequired = schema.ComputedOptional
	default:
		return nil, fmt.Errorf(
			"invalid value for computed_optional_required: %s",
			override.ComputedOptionalRequired,
		)
	}

	return a, nil
}

func (a *ResourceListNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (ResourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *ResourceListNestedAttribute) ToSpec() resource.Attribute {
	a.ListNestedAttribute.NestedObject = resource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return resource.Attribute{
		Name:       util.TerraformIdentifier(a.Name),
		ListNested: &a.ListNestedAttribute,
	}
}

type DataSourceListNestedAttribute struct {
	datasource.ListNestedAttribute

	Name         string
	NestedObject DataSourceNestedAttributeObject
}

func (a *DataSourceListNestedAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceListNestedAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	listNestedAttribute, ok := mergeAttribute.(*DataSourceListNestedAttribute)
	// TODO: return error if types don't match?
	if !ok {
		return a, nil
	}

	if a.Description == nil || *a.Description == "" {
		a.Description = listNestedAttribute.Description
	}
	a.NestedObject.Attributes, _ = a.NestedObject.Attributes.Merge(listNestedAttribute.NestedObject.Attributes)

	return a, nil
}

func (a *DataSourceListNestedAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceListNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (DataSourceAttribute, error) {
	var err error
	a.NestedObject.Attributes, err = a.NestedObject.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *DataSourceListNestedAttribute) ToSpec() datasource.Attribute {
	a.ListNestedAttribute.NestedObject = datasource.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return datasource.Attribute{
		Name:       util.TerraformIdentifier(a.Name),
		ListNested: &a.ListNestedAttribute,
	}
}

type ProviderListNestedAttribute struct {
	provider.ListNestedAttribute

	Name         string
	NestedObject ProviderNestedAttributeObject
}

func (a *ProviderListNestedAttribute) ToSpec() provider.Attribute {
	a.ListNestedAttribute.NestedObject = provider.NestedAttributeObject{
		Attributes: a.NestedObject.Attributes.ToSpec(),
	}

	return provider.Attribute{
		Name:       util.TerraformIdentifier(a.Name),
		ListNested: &a.ListNestedAttribute,
	}
}
