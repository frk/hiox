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

type WriteError struct {
	Err error
}

func (e WriteError) Error() string {
	return e.Err.Error()
}

type ReadError struct {
	Err error
}

func (e ReadError) Error() string {
	return e.Err.Error()
}

type NoTemplateError struct {
	Name string
}

func (e NoTemplateError) Error() string {
	return fmt.Sprintf("hxio: template %q not found", e.Name)
}
