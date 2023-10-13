// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import (
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/datamodel/low"
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

// SchemaErrorFromNode returns a new SchemaError error struct that has no path information, using a schema node to get the line number.
func SchemaErrorFromNode[T any](err error, ref low.NodeReference[T]) *SchemaError {
	lineNumber := 0
	if ref.ValueNode != nil {
		lineNumber = ref.ValueNode.Line
	}

	return &SchemaError{
		err:        err,
		path:       make([]string, 0),
		lineNumber: lineNumber,
	}
}

// SchemaErrorFromProxy returns a new SchemaError error struct that has no path information, using a schema proxy to get the line number.
func SchemaErrorFromProxy(err error, proxy *base.SchemaProxy) *SchemaError {
	lineNumber := 0
	lowProxy := proxy.GoLow()
	if lowProxy != nil && lowProxy.GetValueNode() != nil {
		lineNumber = lowProxy.GetValueNode().Line
	}

	return &SchemaError{
		err:        err,
		path:       make([]string, 0),
		lineNumber: lineNumber,
	}
}
