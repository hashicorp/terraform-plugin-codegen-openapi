// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"gopkg.in/yaml.v3"
)

// SchemaError contains additional details about an error that occurred when processing an OpenAPI schema,
// such as the line number of the invalid schema or nested path information.
type SchemaError struct {
	err        error
	path       []string
	lineNumber int
}

// Error implements the error interface by returning the original error string
func (e *SchemaError) Error() string {
	return e.err.Error()
}

// NestedSchemaError creates a new SchemaError, appending the parent name to the path. This allows a parent
// OpenAPI schema to preserve the error and line number from a child schema, while creating a path name that is an absolute reference.
//
// If no line number exists for the child schema, the parent schema line number will be added.
func (e *SchemaError) NestedSchemaError(parentName string, lineNumber int) *SchemaError {
	newErr := &SchemaError{
		err:        e.err,
		path:       append([]string{parentName}, e.path...),
		lineNumber: e.lineNumber,
	}

	if newErr.lineNumber == 0 {
		newErr.lineNumber = lineNumber
	}

	return newErr
}

// Path returns an absolute reference to the schema where the error occurred.
func (e *SchemaError) Path() string {
	return strings.Join(e.path, ".")
}

// LineNumber returns the line number closest to the schema where the error occurred.
func (e *SchemaError) LineNumber() int {
	return e.lineNumber
}

// SchemaErrorFromProperty returns a new SchemaError error struct
func SchemaErrorFromProperty(err error, name string, lineNumber int) *SchemaError {
	return &SchemaError{
		err:        err,
		path:       []string{name},
		lineNumber: lineNumber,
	}
}

type NodeType int

const (
	None NodeType = iota
	Type
	AdditionalProperties
	Items
	AllOf
	AnyOf
	OneOf
)

// SchemaErrorFromNode returns a new SchemaError error struct that has no path information, using a schema node to get the line number if available.
func SchemaErrorFromNode(err error, schema *base.Schema, nodeType NodeType) *SchemaError {
	// If there is no low information, then we can't retrieve any line numbers
	low := schema.GoLow()
	if low == nil {
		return emptySchemaError(err)
	}

	var valueNode *yaml.Node
	switch nodeType {
	case Type:
		valueNode = low.Type.ValueNode
	case AdditionalProperties:
		valueNode = low.AdditionalProperties.ValueNode
	case Items:
		valueNode = low.Items.ValueNode
	case AllOf:
		valueNode = low.AllOf.ValueNode
	case AnyOf:
		valueNode = low.AnyOf.ValueNode
	case OneOf:
		valueNode = low.OneOf.ValueNode
	}

	lineNumber := 0
	if valueNode != nil {
		lineNumber = valueNode.Line
	}

	return &SchemaError{
		err:        err,
		path:       make([]string, 0),
		lineNumber: lineNumber,
	}
}

// SchemaErrorFromProxy returns a new SchemaError error struct that has no path information, using a schema proxy to get the line number.
func SchemaErrorFromProxy(err error, proxy *base.SchemaProxy) *SchemaError {
	// If there is no low information, then we can't retrieve any line numbers
	lowProxy := proxy.GoLow()
	if lowProxy == nil || lowProxy.GetValueNode() == nil {
		return emptySchemaError(err)
	}

	return &SchemaError{
		err:        err,
		path:       make([]string, 0),
		lineNumber: lowProxy.GetValueNode().Line,
	}
}

// emptySchemaError will return a simple SchemaError struct that contains no additional OAS information about the error
func emptySchemaError(err error) *SchemaError {
	return &SchemaError{
		err:  err,
		path: make([]string, 0),
	}
}
