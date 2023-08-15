// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Implementations of the Explorer interface relate OpenAPIv3 operations to a set of Terraform Provider resource/data source actions (CRUD)
//   - https://spec.openapis.org/oas/latest.html#operation-object
type Explorer interface {
	FindProvider() (Provider, error)
	FindResources() (map[string]Resource, error)
	FindDataSources() (map[string]DataSource, error)
}

type Resource struct {
	CreateOp         *high.Operation
	ReadOp           *high.Operation
	UpdateOp         *high.Operation
	DeleteOp         *high.Operation
	ParameterMatches map[string]string
}

type DataSource struct {
	ReadOp           *high.Operation
	ParameterMatches map[string]string
}

type Provider struct {
	Name        string
	SchemaProxy *base.SchemaProxy
}
