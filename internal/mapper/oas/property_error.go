package oas

import "strings"

type PropertyError struct {
	err        error
	path       []string
	lineNumber int
}

func (e *PropertyError) Error() string {
	return e.err.Error()
}

func (e *PropertyError) NestedPropertyError(name string, lineNumber int) *PropertyError {
	propErr := &PropertyError{
		err:        e.err,
		path:       append([]string{name}, e.path...),
		lineNumber: e.lineNumber,
	}

	// We want to keep the deepest nested line number if it exists
	if e.lineNumber == 0 {
		e.lineNumber = lineNumber
	}

	return propErr
}

func (e *PropertyError) Path() string {
	return strings.Join(e.path, ".")
}

func (e *PropertyError) LineNumber() int {
	return e.lineNumber
}

func NewPropertyError(err error, name string, lineNumber int) *PropertyError {
	return &PropertyError{
		err:        err,
		path:       []string{name},
		lineNumber: lineNumber,
	}
}

func EmptyPropertyError(err error) *PropertyError {
	return &PropertyError{
		err:  err,
		path: make([]string, 0),
	}
}
