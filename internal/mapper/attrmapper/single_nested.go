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

func (a *ResourceSingleNestedAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
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

func (a *ResourceSingleNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (ResourceAttribute, error) {
	var err error
	a.Attributes, err = a.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *ResourceSingleNestedAttribute) ToSpec() resource.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return resource.Attribute{
		Name:         util.TerraformIdentifier(a.Name),
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

func (a *DataSourceSingleNestedAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
	a.Description = &override.Description

	return a, nil
}

func (a *DataSourceSingleNestedAttribute) ApplyNestedOverride(path []string, override explorer.Override) (DataSourceAttribute, error) {
	var err error
	a.Attributes, err = a.Attributes.ApplyOverride(path, override)

	return a, err
}

func (a *DataSourceSingleNestedAttribute) ToSpec() datasource.Attribute {
	a.SingleNestedAttribute.Attributes = a.Attributes.ToSpec()

	return datasource.Attribute{
		Name:         util.TerraformIdentifier(a.Name),
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
		Name:         util.TerraformIdentifier(a.Name),
		SingleNested: &a.SingleNestedAttribute,
	}
}
