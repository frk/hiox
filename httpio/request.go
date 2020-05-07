package httpio

import (
	"context"
	"net/http"
	"net/url"

	"github.com/frk/route"
)

// HeaderReader interface can be implemented to read an incoming request's header.
type HeaderReader interface {
	ReadHeader(header http.Header) error
}

// QueryReader interface can be implemented to read the query parameters
// from an incoming request's url.
type QueryReader interface {
	ReadQuery(query url.Values) error
}

// PathReader interface can be implemented to read the route.Params from an
// incoming request's url path.
type PathReader interface {
	ReadPath(params route.Params) error
}

// BodyReader interface can be implemented to read the body of an incoming request.
type BodyReader interface {
	ReadBody(r *http.Request) error
}

// RequestReader is intended to be embedded as part of a Handler implementation
// to read the various data from the incoming HTTP request.
type RequestReader struct {
	// If set, will read the header from the incoming request.
	Header HeaderReader
	// If set, will read the query parameters from the incoming request's url.
	Query QueryReader
	// If set, will read the route.Params from the path of incoming request's url.
	Path PathReader
	// If set, will read the body from the incoming request.
	Body BodyReader

	// the original request, set by ReadRequest
	r *http.Request
}

// RequestReader implements the ReadRequest method of the Handler interface.
func (rr *RequestReader) ReadRequest(r *http.Request, c context.Context) error {
	rr.r = r

	if rr.Header != nil {
		header := r.Header
		if err := rr.Header.ReadHeader(header); err != nil {
			return err
		}
	}

	if rr.Query != nil {
		query := r.URL.Query()
		if err := rr.Query.ReadQuery(query); err != nil {
			return err
		}
	}

	if rr.Path != nil {
		params := route.GetParams(c)
		if err := rr.Path.ReadPath(params); err != nil {
			return err
		}
	}

	if rr.Body != nil {
		return rr.Body.ReadBody(r)
	}

	return nil
}

// GetContext is a convenience method that returns the underlying http request's context value.
func (rr *RequestReader) GetContext() context.Context {
	return rr.r.Context()
}

// GetRequest is a convenience method that returns the underlying http request.
func (rr *RequestReader) GetRequest() *http.Request {
	return rr.r
}
