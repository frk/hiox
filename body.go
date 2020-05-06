package httpcrud

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/frk/form"
)

// The RequestDump type implements the BodyReader interface.
type RequestDump struct {
	// The pointer to which to set the request dump.
	Val *[]byte
	// Indicates whether to dump the request's body as well.
	Body bool
}

// ReadBody implements the BodyReader interface by dumping the request's
// payload and setting the reciever's Val pointer to the result.
func (d RequestDump) ReadBody(r *http.Request) error {
	dump, err := httputil.DumpRequest(r, d.Body)
	if err != nil {
		return ReadError{err}
	}
	*d.Val = dump
	return nil
}

// The JSON type implements both the BodyWriter and the BodyReader interfaces.
type JSON struct {
	// The value to be json encoded and sent in an HTTP response body or
	// a pointer to the value to be json decoded from an HTTP request's body.
	Val interface{}
}

// ReadBody implements the BodyReader interface by decoding the request's
// json body into the reciever's Val field.
func (j JSON) ReadBody(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(j.Val); err != nil {
		return ReadError{err}
	}
	return nil
}

// WriteInit is a noop, required only to satisfy the BodyWriter interface.
func (JSON) WriteInit(_ http.ResponseWriter) error {
	return nil
}

const contentTypeJSON = "application/json; charset=utf-8"

// WriteBody implements the BodyWriter interface by json encoding the
// receiver's Val field and sending the result in the response's body.
func (j JSON) WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(j.Val); err != nil {
		return WriteError{err}
	}
	return nil
}

// The XML type implements both the BodyWriter and the BodyReader interfaces.
type XML struct {
	// The value to be xml encoded and sent in an HTTP response body or
	// a pointer to the value to be xml decoded from an HTTP request's body.
	Val interface{}
}

// ReadBody implements the BodyReader interface by decoding the request's
// xml body into the reciever's Val field.
func (x XML) ReadBody(r *http.Request) error {
	if err := xml.NewDecoder(r.Body).Decode(x.Val); err != nil {
		return ReadError{err}
	}
	return nil
}

// WriteInit is a noop, required only to satisfy the BodyWriter interface.
func (XML) WriteInit(_ http.ResponseWriter) error {
	return nil
}

const contentTypeXML = "application/xml; charset=utf-8"

// WriteBody implements the BodyWriter interface by xml encoding the
// receiver's Val field and sending the result in the response's body.
func (x XML) WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error {
	w.Header().Set("Content-Type", contentTypeXML)
	w.WriteHeader(statusCode)
	if err := xml.NewEncoder(w).Encode(x.Val); err != nil {
		return WriteError{err}
	}
	return nil
}

// The Form type implements both the BodyWriter and the BodyReader interfaces.
type Form struct {
	// The value to be url encoded and sent in an HTTP response body or
	// a pointer to the value to be url decoded from an HTTP request's body.
	Val interface{}
}

// ReadBody implements the BodyReader interface by decoding the request's
// form body into the reciever's Val field.
func (f Form) ReadBody(r *http.Request) error {
	if err := form.NewDecoder(r.Body).Decode(f.Val); err != nil {
		return ReadError{err}
	}
	return nil
}

// WriteInit is a noop, required only to satisfy the BodyWriter interface.
func (Form) WriteInit(_ http.ResponseWriter) error {
	return nil
}

const contentTypeForm = "application/x-www-form-urlencoded; charset=utf-8"

// WriteBody implements the BodyWriter interface by form encoding the
// receiver's Val field and sending the result in the response's body.
func (f Form) WriteBody(w http.ResponseWriter, r *http.Request, statusCode int) error {
	w.Header().Set("Content-Type", contentTypeForm)
	w.WriteHeader(statusCode)
	if err := form.NewEncoder(w).Encode(f.Val); err != nil {
		return WriteError{err}
	}
	return nil
}

// The Text type implements the BodyWriter interface.
type Text struct {
	// The value to be sent in an HTTP response body.
	Val string
}

// WriteInit is a noop, required only to satisfy the BodyWriter interface.
func (Text) WriteInit(_ http.ResponseWriter) error {
	return nil
}

const contentTypeText = "text/plain"

// WriteBody implements the BodyWriter interface by writing the Val field's
// contents to the response's body.
func (t Text) WriteBody(w http.ResponseWriter, r *http.Request, code int) error {
	w.Header().Set("Content-Type", contentTypeText)
	w.WriteHeader(code)
	_, err := w.Write([]byte(t.Val))
	return err
}

// The HTML type implements the BodyWriter interface.
type HTML struct {
	// The name (map key) of the template as registered with RegisterTemplates.
	Name string
	// Data to be passed to the template's Execute method.
	Data interface{}
}

// WriteInit is a noop, required only to satisfy the BodyWriter interface.
func (HTML) WriteInit(_ http.ResponseWriter) error {
	return nil
}

const contentTypeHTML = "text/html; charset=utf-8"

// WriteBody implements the BodyWriter interface by executing an html
// template with the specified Name, passing it the provided Data, and
// then sending the result as the response's body.
func (h HTML) WriteBody(w http.ResponseWriter, r *http.Request, code int) error {
	if t, ok := templateMap[h.Name]; ok {
		w.Header().Set("Content-Type", contentTypeHTML)
		w.WriteHeader(code)
		return t.Execute(w, h.Data)
	}
	return NoTemplateError{h.Name}
}

var registerHTMLTemplatesOnce sync.Once
var templateMap map[string]*template.Template

// RegisterHTMLTemplatesOnce registers the given templates to be used by the HTML BodyWriter implementation.
// The RegisterHTMLTemplatesOnce function is intended be called once at program start up.
func RegisterHTMLTemplatesOnce(tmap map[string]*template.Template) {
	registerHTMLTemplatesOnce.Do(func() {
		templateMap = make(map[string]*template.Template, len(tmap))
		for key, val := range tmap {
			templateMap[key] = val
		}
	})
}

// The Redirect type implements the BodyWriter interface.
type Redirect struct {
	// The URL to which to redirect.
	URL string
	// The HTTP redirect status code.
	StatusCode int
}

// WriteInit is a noop, required only to satisfy the BodyWriter interface.
func (re Redirect) WriteInit(_ http.ResponseWriter) error {
	return nil
}

// WriteBody implements the BodyWriter interface by issuing a redirect
// to the specified URL with the specified StatusCode.
func (re Redirect) WriteBody(w http.ResponseWriter, r *http.Request, _ int) error {
	http.Redirect(w, r, re.URL, re.StatusCode)
	return nil
}

// The CSV type implements the BodyWriter interface by using a StreamWriter.
type CSV struct {
	Stream StreamWriter
}

// WriteInit opens the underlying StreamWriter.
func (c CSV) WriteInit(w http.ResponseWriter) error {
	c.Stream.Open(w)
	return nil
}

// WriteBody flushes the underlying StreamWriter.
func (c CSV) WriteBody(_ http.ResponseWriter, _ *http.Request, _ int) error {
	if err := c.Stream.Flush(); err != nil {
		return WriteError{err}
	}
	return nil
}

// StreamWriter interface can be used by BodyWriters
type StreamWriter interface {
	// Open, intended to be invoked by (BodyWriter).WriteInit, should prepare
	// the stream for writing using the given http.ResponseWriter as the medium
	// for transferring the data.
	Open(w http.ResponseWriter)
	// Flush, intended to be invoked by (BodyWriter).WriteBody, should write
	// any buffered data to the underlying http.ResponseWriter and return an
	// error if the flush fails.
	Flush() error
}

// CSVWriter implements the StreamWriter interface and is intended to be
// embedded by user-defined structs that want to implement custom csv writers.
type CSVWriter struct {
	// The list of field names.
	// SHOULD be set prior to the first call to WriteRow.
	Header []string
	// The file name to be used for the Content-Disposition header.
	// SHOULD be set prior to the first call to WriteRow.
	FileName string
	// The HTTP status code to be sent with the response.
	// SHOULD be set prior to the first call to WriteRow.
	//
	// Set by Open to 200 (Status OK) as default. If the embedding code needs
	// to send a different value it should implement its own Open method, invoke
	// the embedded Open method to set up the writer and *then* set the StatusCode
	// field to the desired value.
	StatusCode int

	// The http.ResponseWriter target into which to the csv data will be written.
	rw http.ResponseWriter
	// The csv.Writer used to encode and write the given data as csv.
	csv *csv.Writer
	// The write func is used by WriteRow to indirectly invoke
	// the receiver's write1 and write2 methods.
	write func([]string) error
}

// Open prepares the CSVWriter instance using the given http.ResponseWriter
// as the target into which the csv data will be written. Open will be invoked
// by (CSV).WriteInit indirectly through the StreamWriter interface, if, however,
// the embedding code overrides this method by providing its own impelmentation
// then that implementation MUST invoke this method directly.
func (w *CSVWriter) Open(rw http.ResponseWriter) {
	w.rw = rw
	w.write = w.write1
	w.StatusCode = http.StatusOK
}

// Flush flushes the underlying csv.Writer and returns and error if it fails.
// Flush will be invoked by (CSV).WriteBody indirectly through the StreamWriter
// interface, if, however, the embedding code overrides this method by providing
// its own impelmentation then that implementation MUST invoke this method directly.
func (w *CSVWriter) Flush() error {
	w.csv.Flush()
	return w.csv.Error()
}

// WriteRow writes the given row to the underlying http.ResponseWriter.
// This method SHOULD be invoked directly by the user code.
func (w *CSVWriter) WriteRow(row []string) error {
	return w.write(row)
}

const contentTypeCSV = "text/csv"

// write1 sets the csv specific http headers, initializes the underlying
// csv.Writer and writes the list of field names and the given row to it.
func (w *CSVWriter) write1(row []string) error {
	w.rw.Header().Set("Content-Disposition", "attachment; filename="+w.FileName)
	w.rw.Header().Set("Content-Type", contentTypeCSV)
	w.rw.WriteHeader(w.StatusCode)

	w.csv = csv.NewWriter(w.rw)
	if err := w.csv.Write(w.Header); err != nil {
		return err
	}
	if err := w.csv.Write(row); err != nil {
		return err
	}

	w.write = w.write2
	return nil
}

// write2 delegates to the underlying csv.Writer's Write method.
func (w *CSVWriter) write2(row []string) error {
	return w.csv.Write(row)
}
