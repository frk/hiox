package hxio

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/frk/form"
)

type JSON struct {
	Val interface{}
}

// ReadBody implements the BodyReader interface.
func (js JSON) ReadBody(r *http.Request) error {
	if js.Val == nil {
		return nil
	}

	if err := json.NewDecoder(r.Body).Decode(js.Val); err != nil {
		return ReadError{err}
	}
	return nil
}

// WriteInit is a noop. Method is required to implement BodyWriter interface.
func (js JSON) WriteInit(_ http.ResponseWriter) error {
	return nil
}

// WriteBody implements part of the BodyWriter interface.
func (js JSON) WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error {
	const contentType = "application/json; charset=utf-8"

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if js.Val == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(js.Val); err != nil {
		return WriteError{err}
	}
	return nil
}

type XML struct {
	Val interface{}
}

// ReadBody implements the BodyReader interface.
func (x XML) ReadBody(r *http.Request) error {
	if x.Val == nil {
		return nil
	}

	if err := xml.NewDecoder(r.Body).Decode(x.Val); err != nil {
		return ReadError{err}
	}
	return nil
}

// WriteInit is a noop. Method is required to implement BodyWriter interface.
func (x XML) WriteInit(_ http.ResponseWriter) error {
	return nil
}

// WriteBody implements part of the BodyWriter interface.
func (x XML) WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error {
	const contentType = "application/xml; charset=utf-8"

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if x.Val == nil {
		return nil
	}

	if err := xml.NewEncoder(w).Encode(x.Val); err != nil {
		return WriteError{err}
	}
	return nil
}

type Form struct {
	Val interface{}
}

// ReadBody implements the BodyReader interface.
func (f Form) ReadBody(r *http.Request) error {
	if f.Val == nil {
		return nil
	}

	if err := form.NewDecoder(r.Body).Decode(f.Val); err != nil {
		return ReadError{err}
	}
	return nil
}

// WriteInit is a noop. Method is required to implement BodyWriter interface.
func (f Form) WriteInit(_ http.ResponseWriter) error {
	return nil
}

// WriteBody implements part of the BodyWriter interface.
func (f Form) WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error {
	const contentType = "application/x-www-form-urlencoded; charset=utf-8"

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if f.Val == nil {
		return nil
	}

	if err := form.NewEncoder(w).Encode(f.Val); err != nil {
		return WriteError{err}
	}
	return nil
}

// TODO
// - HTML (by template name)
// - HTMLTemplate (with template instance)
// - Redirect
// - CSV
// - CSVStream
