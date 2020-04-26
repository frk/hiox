package hxio

import (
	"context"
	"net/http"
)

type Handler interface {
	AuthCheck(r *http.Request, c context.Context) error
	ReadRequest(r *http.Request, c context.Context) error
	InitResponse(w http.ResponseWriter) error

	BeforeValidate() error
	Validate() error
	AfterValidate() error
	BeforeExecute() error
	Execute() error
	AfterExecute() error

	WriteResponse(w http.ResponseWriter, r *http.Request) error

	Done(e error) error
}

type HandlerInitializer interface {
	Init() Handler
}

// Done can be returned by any handler method to indicate that the
// execution should skip to, and invoke, the handler's Done method
// without calling any of the methods in between.
var Done done

type done struct{}

// implements the error interface.
func (done) Error() string { return `hxio_sigdone` }

type handlerexecer struct {
	init HandlerInitializer
}

// serve initializes and executes the handler.
func (x *handlerexecer) serve(w http.ResponseWriter, r *http.Request, c context.Context) error {
	h := x.init.Init()
	if err := x.exec(h, w, r, c); err != nil && err != Done {
		return h.Done(err)
	}
	return h.Done(nil)
}

// exec invokes the given handler's methods in the pre-defined order, if any of
// the methods return an error exec will exit immediately and return that error,
// leaving the rest of the handler's methods untouched.
func (x *handlerexecer) exec(h Handler, w http.ResponseWriter, r *http.Request, c context.Context) error {
	if err := h.AuthCheck(r, c); err != nil {
		return err
	}
	if err := h.ReadRequest(r, c); err != nil {
		return err
	}
	if err := h.InitResponse(w); err != nil {
		return err
	}

	if err := h.BeforeValidate(); err != nil {
		return err
	}
	if err := h.Validate(); err != nil {
		return err
	}
	if err := h.AfterValidate(); err != nil {
		return err
	}
	if err := h.BeforeExecute(); err != nil {
		return err
	}
	if err := h.Execute(); err != nil {
		return err
	}
	if err := h.AfterExecute(); err != nil {
		return err
	}

	return h.WriteResponse(w, r)
}
