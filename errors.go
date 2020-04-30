package hxio

import (
	"fmt"
)

// The IsDone "signal" can be returned by any Action method to indicate that
// the execution should skip to, and invoke, the Action's Done method without
// calling any of its other methods in-between.
var IsDone done

type done struct{}

// implements the error interface.
func (done) Error() string { return `hxio:sigdone` }

// WriteError represents an error returned by a BodyWriter.
type WriteError struct {
	// The original error.
	Err error
}

func (e WriteError) Error() string {
	return e.Err.Error()
}

// ReadError represents an error returned by a BodyReader.
type ReadError struct {
	// The original error.
	Err error
}

func (e ReadError) Error() string {
	return e.Err.Error()
}

// NoTemplateError is returned when no template with the given name was registered.
type NoTemplateError struct {
	// The provided template name.
	Name string
}

func (e NoTemplateError) Error() string {
	return fmt.Sprintf("hxio: template %q not found", e.Name)
}
