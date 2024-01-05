// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/explorer"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
)

type ResourceAttribute interface {
	GetName() string
	Merge(ResourceAttribute) (ResourceAttribute, error)
	ApplyOverride(explorer.Override) (ResourceAttribute, error)
	ToSpec() resource.Attribute
}

type ResourceNestedAttribute interface {
	ApplyNestedOverride([]string, explorer.Override) (ResourceAttribute, error)
	NestedMerge([]string, ResourceAttribute, schema.ComputedOptionalRequired) (ResourceAttribute, error)
}

type ResourceAttributes []ResourceAttribute

func (targetSlice ResourceAttributes) Merge(mergeSlices ...ResourceAttributes) (ResourceAttributes, error) {
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

func (targetSlice ResourceAttributes) MergeAttribute(path []string, attribute ResourceAttribute, intermediateComputability schema.ComputedOptionalRequired) (ResourceAttributes, error) {
	var errResult error
	if len(path) == 0 {
		return targetSlice, errResult
	}
	for i, target := range targetSlice {
		if target.GetName() == path[0] {

			if len(path) > 1 {
				nestedTarget, ok := target.(ResourceNestedAttribute)
				if !ok {
					// TODO: error? there is a nested override for an attribute that is not a nested type
					break
				}

				// The attribute we need to override is deeper nested, move up
				nextPath := path[1:]

				overriddenTarget, err := nestedTarget.NestedMerge(nextPath, attribute, intermediateComputability)
				errResult = errors.Join(errResult, err)

				targetSlice[i] = overriddenTarget

			} else {
				// No more path to traverse, apply merge, bidirectional
				overriddenTarget, err := attribute.Merge(target)
				errResult = errors.Join(errResult, err)
				overriddenTarget, err = overriddenTarget.Merge(attribute)
				errResult = errors.Join(errResult, err)

				targetSlice[i] = overriddenTarget
			}

			break
		}
	}

	return targetSlice, errResult
}

func (attributes ResourceAttributes) ToSpec() []resource.Attribute {
	specAttributes := make([]resource.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}

func (attributes ResourceAttributes) ApplyOverrides(overrideMap map[string]explorer.Override) (ResourceAttributes, error) {
	var errResult error
	for key, override := range overrideMap {
		var err error
		attributes, err = attributes.ApplyOverride(strings.Split(key, "."), override)
		errResult = errors.Join(errResult, err)
	}

	return attributes, errResult
}

func (attributes ResourceAttributes) ApplyOverride(path []string, override explorer.Override) (ResourceAttributes, error) {
	var errResult error
	if len(path) == 0 {
		return attributes, errResult
	}
	for i, attribute := range attributes {
		if attribute.GetName() == path[0] {

			if len(path) > 1 {
				nestedAttribute, ok := attribute.(ResourceNestedAttribute)
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
