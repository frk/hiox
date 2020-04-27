package hxio

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/form"
)

type strReadCloser struct{ *strings.Reader }

func (strReadCloser) Close() error { return nil }

type testBody struct {
	XMLName xml.Name `json:"-" xml:"data" form:"-"`
	Foo     string   `json:"foo" xml:"foo" form:"foo"`
	Bar     float64  `json:"bar" xml:"bar" form:"bar"`
	Baz     bool     `json:"baz" xml:"baz" form:"baz"`
}

func TestJSON_ReadBody(t *testing.T) {
	plaintext := `plaintext`
	syntaxError := json.Unmarshal([]byte(plaintext), &testBody{})

	tests := []struct {
		name  string
		body  string
		input interface{}
		want  interface{}
		err   error
	}{{
		name:  "should fail when empty body",
		body:  ``,
		input: &testBody{},
		want:  &testBody{},
		err:   ReadError{io.EOF},
	}, {
		name:  "should fail when non-json body",
		body:  plaintext,
		input: &testBody{},
		want:  &testBody{},
		err:   ReadError{syntaxError},
	}, {
		name:  "should decode json body into Value",
		body:  `{"foo":"test","bar":0.004,"baz":true}`,
		input: &testBody{},
		want:  &testBody{Foo: "test", Bar: 0.004, Baz: true},
	}, {
		name:  "should not decode if no Value was provided",
		body:  `{"foo":"test","bar":0.004,"baz":true}`,
		input: nil,
		want:  nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Body: strReadCloser{strings.NewReader(tt.body)}}

			err := (&JSON{tt.input}).ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.input, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestJSON_WriteBody(t *testing.T) {
	const contentType = "application/json; charset=utf-8"

	tests := []struct {
		name   string
		data   *testBody
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write null",
		code:   204,
		want:   "null\n",
		header: http.Header{"Content-Type": {contentType}},
	}, {
		name:   "write data",
		code:   200,
		data:   &testBody{Foo: "test", Bar: 0.004, Baz: true},
		want:   `{"foo":"test","bar":0.004,"baz":true}` + "\n",
		header: http.Header{"Content-Type": {contentType}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := (&JSON{tt.data}).WriteBody(w, nil, tt.code)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Code, tt.code); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Body.String(), tt.want); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Result().Header, tt.header); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestXML_ReadBody(t *testing.T) {
	plaintext := `plaintext`
	syntaxError := xml.Unmarshal([]byte(plaintext), &testBody{})

	tests := []struct {
		name  string
		body  string
		input interface{}
		want  interface{}
		err   error
	}{{
		name:  "should fail when empty body",
		body:  ``,
		input: &testBody{},
		want:  &testBody{},
		err:   ReadError{io.EOF},
	}, {
		name:  "should fail when non-xml body",
		body:  plaintext,
		input: &testBody{},
		want:  &testBody{},
		err:   ReadError{syntaxError},
	}, {
		name:  "should decode xml body",
		body:  `<data><foo>test</foo><bar>0.004</bar><baz>true</baz></data>`,
		input: &testBody{},
		want:  &testBody{XMLName: xml.Name{Local: "data"}, Foo: "test", Bar: 0.004, Baz: true},
	}, {
		name:  "should not decode xml body if no Value was provided",
		body:  `<data><foo>test</foo><bar>0.004</bar><baz>true</baz></data>`,
		input: nil,
		want:  nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Body: strReadCloser{strings.NewReader(tt.body)}}

			err := (&XML{tt.input}).ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.input, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestXML_WriteBody(t *testing.T) {
	const contentType = "application/xml; charset=utf-8"

	tests := []struct {
		name   string
		data   *testBody
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write empty string",
		code:   204,
		want:   "",
		header: http.Header{"Content-Type": {contentType}},
	}, {
		name:   "write data",
		code:   200,
		data:   &testBody{Foo: "test", Bar: 0.004, Baz: true},
		want:   `<data><foo>test</foo><bar>0.004</bar><baz>true</baz></data>`,
		header: http.Header{"Content-Type": {contentType}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := (&XML{tt.data}).WriteBody(w, nil, tt.code)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Code, tt.code); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Body.String(), tt.want); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Result().Header, tt.header); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestForm_ReadBody(t *testing.T) {
	badtypetext := `baz=abc def`
	valueError := form.Unmarshal([]byte(badtypetext), &testBody{})

	tests := []struct {
		name  string
		body  string
		input interface{}
		want  interface{}
		err   error
	}{{
		name: "ignore empty body",
		body: ``,
		// input: &testBody{},
		// want: &testBody{},
	}, {
		name:  "should fail with incompatible types",
		body:  badtypetext,
		input: &testBody{},
		want:  &testBody{},
		err:   ReadError{valueError},
	}, {
		name:  "should decode form body",
		body:  `foo=test&bar=0.004&baz=true`,
		input: &testBody{},
		want:  &testBody{Foo: "test", Bar: 0.004, Baz: true},
	}, {
		name:  "should not decode form body if no dest input was set",
		body:  `foo=test&bar=0.004&baz=true`,
		input: nil,
		want:  nil,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Body: strReadCloser{strings.NewReader(tt.body)}}

			err := (&Form{Val: tt.input}).ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.input, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestForm_WriteBody(t *testing.T) {
	const contentType = "application/x-www-form-urlencoded; charset=utf-8"

	tests := []struct {
		name   string
		data   *testBody
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "fail with bad-value",
		code:   204,
		want:   "",
		header: http.Header{"Content-Type": {contentType}},
	}, {
		name:   "write data",
		code:   200,
		data:   &testBody{Foo: "test", Bar: 0.004, Baz: true},
		want:   `foo=test&bar=0.004&baz=true`,
		header: http.Header{"Content-Type": {contentType}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := (&Form{Val: tt.data}).WriteBody(w, nil, tt.code)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e, err)
			}
			if e := compare.Compare(w.Code, tt.code); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Body.String(), tt.want); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(w.Result().Header, tt.header); e != nil {
				t.Error(e)
			}
		})
	}
}
