package httpio

import (
	"fmt"
)

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
	return fmt.Sprintf("httpcrud/httpio: template %q not found", e.Name)
}
