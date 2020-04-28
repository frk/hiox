package hxio

import (
	"fmt"
)

type WriteError struct {
	Err error
}

func (we WriteError) Error() string {
	return we.Err.Error()
}

type ReadError struct {
	Err error
}

func (re ReadError) Error() string {
	return re.Err.Error()
}

type NoTemplateError struct {
	Name string
}

func (te NoTemplateError) Error() string {
	return fmt.Sprintf("hxio: template %q not found", te.Name)
}
