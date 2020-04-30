package hxio

import (
	"context"
	"net/http"
)

type HandlerInitializerFactory interface {
	//
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

type handler struct{ action }

func (handler) AuthCheck(_ *http.Request, _ context.Context) error         { return nil }
func (handler) ReadRequest(_ *http.Request, _ context.Context) error       { return nil }
func (handler) InitResponse(_ http.ResponseWriter) error                   { return nil }
func (handler) WriteResponse(_ http.ResponseWriter, _ *http.Request) error { return nil }

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
