// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package explorer

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Explorer implements methods that relate OpenAPI operations to a set of Terraform Provider resource/data source actions (CRUD)
type Explorer interface {
	FindProvider() (Provider, error)
	FindResources() (map[string]Resource, error)
	FindDataSources() (map[string]DataSource, error)
}

// Resource contains CRUD operations and schema options for configuration.
type Resource struct {
	CreateOp         *high.Operation
	ReadOp           *high.Operation
	UpdateOp         *high.Operation
	DeleteOp         *high.Operation
	CommonParameters []*high.Parameter
	SchemaOptions    SchemaOptions
}

// DataSource contains a Read operation and schema options for configuration.
type DataSource struct {
	ReadOp           *high.Operation
	CommonParameters []*high.Parameter
	SchemaOptions    SchemaOptions
}

// Provider contains a name and a schema.
type Provider struct {
	Name        string
	SchemaProxy *base.SchemaProxy
	Ignores     []string
}

type SchemaOptions struct {
	Ignores          []string
	AttributeOptions AttributeOptions
}

type AttributeOptions struct {
	Aliases   map[string]string
	Overrides map[string]Override
}

type Override struct {
	Description              string
	ComputedOptionalRequired string
}
