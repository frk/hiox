package hxio

import (
	"context"
	"net/http"
	"net/url"

	"github.com/frk/route"
)

type HeaderReader interface {
	ReadHeader(header http.Header) error
}

type QueryReader interface {
	ReadQuery(query url.Values) error
}

type PathReader interface {
	ReadPath(params route.Params) error
}

type BodyReader interface {
	ReadBody(r *http.Request) error
}

type RequestReader struct {
	Header HeaderReader
	Query  QueryReader
	Path   PathReader
	Body   BodyReader
}

func (rr *RequestReader) ReadRequest(r *http.Request, c context.Context) error {
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
