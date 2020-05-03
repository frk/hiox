package hiox

import (
	"net/http"
)

// HeaderWriter interface can be implemented to write an outgoing response's header.
type HeaderWriter interface {
	WriteHeader(header http.Header)
}

// BodyWriter interface can be implemented to write the body of an outgoing response.
type BodyWriter interface {
	WriteInit(w http.ResponseWriter) error
	WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error
}

// ResponseWriter is intended to be embedded as part of a Handler implementation
// to write the various data of the outgoing HTTP response.
type ResponseWriter struct {
	// If set, will write the header of the outgoing response.
	Header HeaderWriter
	// If set, will write the body of the outgoing response.
	Body BodyWriter
	// If set, will be used as the HTTP status code of the outgoing response.
	Status int
}

// ResponseWriter implements the InitResponse method of the Handler interface.
func (rw *ResponseWriter) InitResponse(w http.ResponseWriter) error {
	if rw.Body != nil {
		return rw.Body.WriteInit(w)
	}
	return nil
}

// ResponseWriter implements the WriteResponse method of the Handler interface.
func (rw *ResponseWriter) WriteResponse(w http.ResponseWriter, r *http.Request) error {
	if rw.Header != nil {
		rw.Header.WriteHeader(w.Header())
	}

	if rw.Body != nil {
		if rw.Status > 0 {
			return rw.Body.WriteBody(w, r, rw.Status)
		}
		return rw.Body.WriteBody(w, r, http.StatusOK) // default to 200
	} else if rw.Status > 0 {
		w.WriteHeader(rw.Status)
	}
	return nil
}
