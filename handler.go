package hxio

import (
	"context"
	"net/http"
)

type HandlerInitializerFactory interface {
	NewHandlerInitializer(v interface{}, path, method string) HandlerInitializer
}

// HandlerInitializer
type HandlerInitializer interface {
	// Init
	Init() Handler
}

// Handler
type Handler interface {
	AuthCheck(r *http.Request, c context.Context) error
	ReadRequest(r *http.Request, c context.Context) error
	InitResponse(w http.ResponseWriter) error
	Action
	WriteResponse(w http.ResponseWriter, r *http.Request) error
}

// handlerExecer
type handlerExecer struct {
	init HandlerInitializer
}

// serve initializes and executes the handler.
func (x *handlerExecer) serve(w http.ResponseWriter, r *http.Request, c context.Context) error {
	h := x.init.Init()
	return x.exec(h, w, r, c)
}

// exec invokes the given handler's methods in the pre-defined order, if any of
// the methods return an error exec will exit immediately and return that error,
// leaving the rest of the handler's methods untouched.
func (x *handlerExecer) exec(h Handler, w http.ResponseWriter, r *http.Request, c context.Context) error {
	if err := h.AuthCheck(r, c); err != nil {
		return err
	}
	if err := h.ReadRequest(r, c); err != nil {
		return err
	}
	if err := h.InitResponse(w); err != nil {
		return err
	}

	if err := ExecuteAction(h); err != nil {
		return err
	}

	return h.WriteResponse(w, r)
}

// HandlerBase is a noop helper type that can be embedded by user defined
// types that are intended to implement the Handler interface but do not
// need to, nor want to, declare every single one of its methods.
type HandlerBase struct{ handlerbase }

// handlerbase is embedded by HandlerBase to artificially increase the depth level
// of the noop methods to reduce the possibility of an "ambiguous selector" issue.
type handlerbase struct{ actionbase }

func (handlerbase) AuthCheck(_ *http.Request, _ context.Context) error         { return nil }
func (handlerbase) ReadRequest(_ *http.Request, _ context.Context) error       { return nil }
func (handlerbase) InitResponse(_ http.ResponseWriter) error                   { return nil }
func (handlerbase) WriteResponse(_ http.ResponseWriter, _ *http.Request) error { return nil }
