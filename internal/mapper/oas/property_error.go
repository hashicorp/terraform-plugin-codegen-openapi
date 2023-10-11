// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas

import "strings"

// PropertyError contains additional details about an error that occurred when processing a specific OpenAPI property,
// such as a line number for the property or nested path information.
type PropertyError struct {
	err        error
	path       []string
	lineNumber int
}

// Error implements the error interface using the original error string
func (e *PropertyError) Error() string {
	return e.err.Error()
}

// NestedPropertyError creates a new PropertyError, appending the parent property name to the path. This allows a parent
// OpenAPI property to preserve the error and line number of a child property, while creating a path name that is an absolute reference
// to the child property.
//
// If no line number exists for the child property, the parent property line number will be added.
func (e *PropertyError) NestedPropertyError(parentPropertyName string, lineNumber int) *PropertyError {
	propErr := &PropertyError{
		err:        e.err,
		path:       append([]string{parentPropertyName}, e.path...),
		lineNumber: e.lineNumber,
	}

	if e.lineNumber == 0 {
		e.lineNumber = lineNumber
	}

	return propErr
}

// Path returns an absolute reference to the property where the error occurred.
func (e *PropertyError) Path() string {
	return strings.Join(e.path, ".")
}

// LineNumber returns the line number closest to the property where the error occurred.
func (e *PropertyError) LineNumber() int {
	return e.lineNumber
}

// NewPropertyError returns a new Property error struct
func NewPropertyError(err error, name string, lineNumber int) *PropertyError {
	return &PropertyError{
		err:        err,
		path:       []string{name},
		lineNumber: lineNumber,
	}
}

// EmptyPropertyError returns a new Property error struct that has no path information. This is useful for errors that occur
// deeply nested in a schema, but may not be immediately adjacent to a property.
func EmptyPropertyError(err error) *PropertyError {
	return &PropertyError{
		err:  err,
		path: make([]string, 0),
	}
}
