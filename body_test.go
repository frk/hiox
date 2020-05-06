package httpcrud

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/form"
)

func TestRequestDump_ReadBody(t *testing.T) {
	//plaintext := `plaintext`
	//syntaxError := json.Unmarshal([]byte(plaintext), &testBody{})

	tests := []struct {
		name   string
		url    string
		method string
		body   string
		dump   RequestDump
		want   []byte
		err    error
	}{{
		name:   "get without body",
		url:    "https://testing.com",
		method: "GET",
		dump:   RequestDump{&[]byte{}, false},
		want:   []byte("GET / HTTP/1.1\r\nHost: testing.com\r\n\r\n"),
	}, {
		name:   "post without body",
		url:    "https://testing.com/a/b/c",
		method: "POST",
		body:   `{"foo":"bar"}`,
		dump:   RequestDump{&[]byte{}, false},
		want:   []byte("POST /a/b/c HTTP/1.1\r\nHost: testing.com\r\n\r\n"),
	}, {
		name:   "post with body",
		url:    "https://testing.com/a/b/c",
		method: "POST",
		body:   `{"foo":"bar"}`,
		dump:   RequestDump{&[]byte{}, true},
		want:   []byte("POST /a/b/c HTTP/1.1\r\nHost: testing.com\r\n\r\n{\"foo\":\"bar\"}"),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			err = tt.dump.ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(string(*tt.dump.Val), string(tt.want)); e != nil {
				fmt.Printf("%q\n", *tt.dump.Val)
				t.Error(e)
			}
		})
	}
}

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
		name string
		body string
		json JSON
		want interface{}
		err  error
	}{{
		name: "should fail when empty body",
		body: ``,
		json: JSON{&testBody{}},
		want: &testBody{},
		err:  ReadError{io.EOF},
	}, {
		name: "should fail when non-json body",
		body: plaintext,
		json: JSON{&testBody{}},
		want: &testBody{},
		err:  ReadError{syntaxError},
	}, {
		name: "should decode json body into Value",
		body: `{"foo":"test","bar":0.004,"baz":true}`,
		json: JSON{&testBody{}},
		want: &testBody{Foo: "test", Bar: 0.004, Baz: true},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Body: strReadCloser{strings.NewReader(tt.body)}}

			err := tt.json.ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.json.Val, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestJSON_WriteBody(t *testing.T) {
	tests := []struct {
		name   string
		json   JSON
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write data",
		code:   200,
		json:   JSON{testBody{Foo: "test", Bar: 0.004, Baz: true}},
		want:   `{"foo":"test","bar":0.004,"baz":true}` + "\n",
		header: http.Header{"Content-Type": {contentTypeJSON}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := tt.json.WriteBody(w, nil, tt.code)
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
		name string
		body string
		xml  XML
		want interface{}
		err  error
	}{{
		name: "should fail when empty body",
		body: ``,
		xml:  XML{&testBody{}},
		want: &testBody{},
		err:  ReadError{io.EOF},
	}, {
		name: "should fail when non-xml body",
		body: plaintext,
		xml:  XML{&testBody{}},
		want: &testBody{},
		err:  ReadError{syntaxError},
	}, {
		name: "should decode xml body",
		body: `<data><foo>test</foo><bar>0.004</bar><baz>true</baz></data>`,
		xml:  XML{&testBody{}},
		want: &testBody{XMLName: xml.Name{Local: "data"}, Foo: "test", Bar: 0.004, Baz: true},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Body: strReadCloser{strings.NewReader(tt.body)}}

			err := tt.xml.ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.xml.Val, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestXML_WriteBody(t *testing.T) {
	tests := []struct {
		name   string
		xml    XML
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write data",
		code:   200,
		xml:    XML{&testBody{Foo: "test", Bar: 0.004, Baz: true}},
		want:   `<data><foo>test</foo><bar>0.004</bar><baz>true</baz></data>`,
		header: http.Header{"Content-Type": {contentTypeXML}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := tt.xml.WriteBody(w, nil, tt.code)
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
		name string
		body string
		form Form
		want interface{}
		err  error
	}{{
		name: "should fail with incompatible types",
		body: badtypetext,
		form: Form{&testBody{}},
		want: &testBody{},
		err:  ReadError{valueError},
	}, {
		name: "should decode form body",
		body: `foo=test&bar=0.004&baz=true`,
		form: Form{&testBody{}},
		want: &testBody{Foo: "test", Bar: 0.004, Baz: true},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Body: strReadCloser{strings.NewReader(tt.body)}}

			err := tt.form.ReadBody(r)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e)
			}
			if e := compare.Compare(tt.form.Val, tt.want); e != nil {
				t.Error(e)
			}
		})
	}
}

func TestForm_WriteBody(t *testing.T) {
	tests := []struct {
		name   string
		form   Form
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write data",
		form:   Form{&testBody{Foo: "test", Bar: 0.004, Baz: true}},
		code:   200,
		want:   `foo=test&bar=0.004&baz=true`,
		header: http.Header{"Content-Type": {contentTypeForm}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := tt.form.WriteBody(w, nil, tt.code)
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

func TestText_WriteBody(t *testing.T) {
	tests := []struct {
		name   string
		text   Text
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write text",
		text:   Text{"hello world"},
		code:   200,
		want:   `hello world`,
		header: http.Header{"Content-Type": {contentTypeText}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := tt.text.WriteBody(w, nil, tt.code)
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

func TestHTML_WriteBody(t *testing.T) {
	t1 := template.Must(template.New("t").Parse(`<html><body>{{ . }}</body></html>`))
	t2 := template.Must(template.New("t").Parse(`<html><body style="background-color:{{ . }}">foo</body></html>`))

	RegisterHTMLTemplatesOnce(map[string]*template.Template{
		"page_a": t1,
		"page_b": t2,
	})

	tests := []struct {
		name   string
		html   HTML
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name:   "write html",
		html:   HTML{"page_a", 93459672},
		code:   200,
		want:   `<html><body>93459672</body></html>`,
		header: http.Header{"Content-Type": {contentTypeHTML}},
	}, {
		name:   "write html 2",
		html:   HTML{"page_a", "Hello!"},
		code:   201,
		want:   `<html><body>Hello!</body></html>`,
		header: http.Header{"Content-Type": {contentTypeHTML}},
	}, {
		name:   "write html 3",
		html:   HTML{"page_b", "yellow"},
		code:   200,
		want:   `<html><body style="background-color:yellow">foo</body></html>`,
		header: http.Header{"Content-Type": {contentTypeHTML}},
	}, {
		name:   "no template error",
		html:   HTML{"page_c", "yellow"},
		header: http.Header{},
		err:    NoTemplateError{"page_c"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			w.Code = 0

			err := tt.html.WriteBody(w, nil, tt.code)
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

func TestRedirect_WriteBody(t *testing.T) {
	tests := []struct {
		name   string
		redir  Redirect
		want   string
		header http.Header
		err    error
	}{{
		name:  "write redirect",
		redir: Redirect{"https://target.com/foobar", 301},
		want:  "<a href=\"https://target.com/foobar\">Moved Permanently</a>.\n\n",
		header: http.Header{"Location": {"https://target.com/foobar"},
			"Content-Type": {contentTypeHTML}},
	}, {
		name:  "write redirect 2",
		redir: Redirect{"/foo-bar", 302},
		want:  "<a href=\"/foo-bar\">Found</a>.\n\n",
		header: http.Header{"Location": {"/foo-bar"},
			"Content-Type": {contentTypeHTML}},
	}, {
		name:  "write redirect 2",
		redir: Redirect{"bar/baz", 303},
		want:  "<a href=\"/foo/bar/baz\">See Other</a>.\n\n",
		header: http.Header{"Location": {"/foo/bar/baz"},
			"Content-Type": {contentTypeHTML}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "https://example.com/foo/", nil)
			w := httptest.NewRecorder()
			w.Code = 0

			err := tt.redir.WriteBody(w, r, 0)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Error(e, err)
			}
			if e := compare.Compare(w.Code, tt.redir.StatusCode); e != nil {
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

func TestCSV_WriteBody(t *testing.T) {
	tests := []struct {
		name   string
		csv    CSV
		writer *CSVWriter
		data   [][]string
		code   int
		want   string
		header http.Header
		err    error
	}{{
		name: "write csv",
		writer: &CSVWriter{
			Header:   []string{"foo", "bar", "baz"},
			FileName: "data.csv",
		},
		data: [][]string{{"abc", "def", "ghi"}, {"123", "456", "789"}},
		code: 200,
		want: "foo,bar,baz\nabc,def,ghi\n123,456,789\n",
		header: http.Header{"Content-Disposition": {"attachment; filename=data.csv"},
			"Content-Type": {contentTypeCSV}},
	}, {
		name: "write csv 2",
		writer: &CSVWriter{
			Header:   []string{"foo", "bar", "baz"},
			FileName: "file.csv",
		},
		data: [][]string{{"123", "456", "789"}, {"abc", "def", "ghi"}},
		code: 201,
		want: "foo,bar,baz\n123,456,789\nabc,def,ghi\n",
		header: http.Header{"Content-Disposition": {"attachment; filename=file.csv"},
			"Content-Type": {contentTypeCSV}},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			w.Code = 0

			tt.csv.Stream = tt.writer
			_ = tt.csv.WriteInit(w)
			tt.writer.StatusCode = tt.code
			for _, row := range tt.data {
				if err := tt.writer.WriteRow(row); err != nil {
					t.Fatal(err)
				}
			}
			err := tt.csv.WriteBody(nil, nil, 0)
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
