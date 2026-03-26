// Copyright IBM Corp. 2023, 2026
// SPDX-License-Identifier: MPL-2.0

package attrmapper

import (
	"github.com/hashicorp/terraform-plugin-codegen-spec/provider"
)

type ProviderAttribute interface {
	ToSpec() provider.Attribute
}

type ProviderAttributes []ProviderAttribute

func (attributes ProviderAttributes) ToSpec() []provider.Attribute {
	specAttributes := make([]provider.Attribute, 0, len(attributes))
	for _, attribute := range attributes {
		specAttributes = append(specAttributes, attribute.ToSpec())
	}

	return specAttributes
}
