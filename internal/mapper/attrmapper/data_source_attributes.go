// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/config"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
)

type DataSourceAttribute interface {
	GetName() string
	Merge(DataSourceAttribute) (DataSourceAttribute, error)
	ApplyOverride(config.Override) (DataSourceAttribute, error)
	ToSpec() datasource.Attribute
}

type DataSourceNestedAttribute interface {
	ApplyNestedOverride([]string, config.Override) (DataSourceAttribute, error)
}

type DataSourceAttributes []DataSourceAttribute

func (targetSlice DataSourceAttributes) Merge(mergeSlices ...DataSourceAttributes) (DataSourceAttributes, error) {
	var errResult error

	for _, mergeSlice := range mergeSlices {
		for _, mergeAttribute := range mergeSlice {
			// As we compare attributes, if we don't find a match, we should add this attribute to the slice after
			isNewAttribute := true

			for i, targetAttribute := range targetSlice {
				if targetAttribute.GetName() == mergeAttribute.GetName() {
					mergedAttribute, err := targetAttribute.Merge(mergeAttribute)
					if err != nil {
						// TODO: consider how best to surface this error
						// Currently, if the merge fails we should just keep the original target attribute for now
						errResult = errors.Join(errResult, err)
					} else {
						targetSlice[i] = mergedAttribute
					}

					isNewAttribute = false
					break
				}
			}

			if isNewAttribute {
				// Add this back to the original slice to avoid adding duplicate attributes from different mergeSlices
				targetSlice = append(targetSlice, mergeAttribute)
			}
		}
	}

	return targetSlice, errResult
}

func (attributes DataSourceAttributes) ToSpec() []datasource.Attribute {
	specAttributes := make([]datasource.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}

func (attributes DataSourceAttributes) ApplyOverrides(overrideMap map[string]config.Override) (DataSourceAttributes, error) {
	var errResult error
	for key, override := range overrideMap {
		var err error
		attributes, err = attributes.ApplyOverride(strings.Split(key, "."), override)
		errResult = errors.Join(errResult, err)
	}

	return attributes, errResult
}

func (attributes DataSourceAttributes) ApplyOverride(path []string, override config.Override) (DataSourceAttributes, error) {
	var errResult error
	if len(path) == 0 {
		return attributes, errResult
	}
	for i, attribute := range attributes {
		if attribute.GetName() == path[0] {

			if len(path) > 1 {
				nestedAttribute, ok := attribute.(DataSourceNestedAttribute)
				if !ok {
					// TODO: error? there is a nested override for an attribute that is not a nested type
					break
				}

				// The attribute we need to override is deeper nested, move up
				nextPath := path[1:]

				overriddenAttribute, err := nestedAttribute.ApplyNestedOverride(nextPath, override)
				errResult = errors.Join(errResult, err)

				attributes[i] = overriddenAttribute

			} else {
				// No more path to traverse, apply override
				overriddenAttribute, err := attribute.ApplyOverride(override)
				errResult = errors.Join(errResult, err)

				attributes[i] = overriddenAttribute
			}

			break
		}
	}

	return attributes, errResult
}
