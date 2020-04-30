package hxio

import (
	"net/http"
)

type HeaderWriter interface {
	WriteHeader(header http.Header)
}

type BodyWriter interface {
	WriteInit(w http.ResponseWriter) error
	WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error
}

type ResponseWriter struct {
	Header HeaderWriter
	Body   BodyWriter
	Status int
}

func (rw *ResponseWriter) InitResponse(w http.ResponseWriter) error {
	if rw.Body != nil {
		return rw.Body.WriteInit(w)
	}
	return nil
}

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
