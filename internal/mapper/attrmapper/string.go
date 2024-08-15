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

type ResourceStringAttribute struct {
	resource.StringAttribute

	Name string
}

func (a *ResourceStringAttribute) GetName() string {
	return a.Name
}

func (a *ResourceStringAttribute) Merge(mergeAttribute ResourceAttribute) (ResourceAttribute, error) {
	stringAttribute, ok := mergeAttribute.(*ResourceStringAttribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = stringAttribute.Description
	}

	return a, nil
}

func (a *ResourceStringAttribute) ApplyOverride(override explorer.Override) (ResourceAttribute, error) {
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

func (a *ResourceStringAttribute) ToSpec() resource.Attribute {
	return resource.Attribute{
		Name:   util.TerraformIdentifier(a.Name),
		String: &a.StringAttribute,
	}
}

type DataSourceStringAttribute struct {
	datasource.StringAttribute

	Name string
}

func (a *DataSourceStringAttribute) GetName() string {
	return a.Name
}

func (a *DataSourceStringAttribute) Merge(mergeAttribute DataSourceAttribute) (DataSourceAttribute, error) {
	stringAttribute, ok := mergeAttribute.(*DataSourceStringAttribute)
	// TODO: return error if types don't match?
	if ok && (a.Description == nil || *a.Description == "") {
		a.Description = stringAttribute.Description
	}

	return a, nil
}

func (a *DataSourceStringAttribute) ApplyOverride(override explorer.Override) (DataSourceAttribute, error) {
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

func (a *DataSourceStringAttribute) ToSpec() datasource.Attribute {
	return datasource.Attribute{
		Name:   util.TerraformIdentifier(a.Name),
		String: &a.StringAttribute,
	}
}

type ProviderStringAttribute struct {
	provider.StringAttribute

	Name string
}

func (a *ProviderStringAttribute) ToSpec() provider.Attribute {
	return provider.Attribute{
		Name:   util.TerraformIdentifier(a.Name),
		String: &a.StringAttribute,
	}
}
