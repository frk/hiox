package hiox

import (
	"context"
	"net/http"
)

// HandlerInitializer is an interface used for initializing new, request
// scoped, instances of the Handler interface.
type HandlerInitializer interface {
	// Init returns a Handler instance (in general a new one).
	//
	// Init will be invoked once for every incoming request and the returned
	// Handler will be used to handle that request. As long as Init is implemented
	// to return a new Handler instance everytime it's invoked it is guaranteed
	// that each request will be handled by their own Handler instance.
	Init() Handler
}

// Handler wraps a set of methods that are executed in sequence to handle an
// incoming HTTP request. Each of the Handler's methods is intended to implement
// a subset of the work that the Handler needs to do to handle the request.
// The kind of the work that a method should perform is indicated by the method's
// name, however this is nonbinding and it is left to the developer's judgement
// to decide which method should do what subset of the Handler's work.
//
// If, during execution, any of the Handler's methods returns an error the execution
// will stop and return that error, leaving all of the subsequent methods uninvoked.
type Handler interface {
	// Authenticate and authorize the incoming request.
	AuthCheck(r *http.Request, c context.Context) error
	// Read the input from the incoming request (headers, url, query, body).
	ReadRequest(r *http.Request, c context.Context) error
	// Prepare the response. In the general case this is unnecessary and
	// most handlers do not need to do anything in this method.
	//
	// However this method becomes useful when one needs to stream data,
	// more specifically it allows the Handler to setup the writer as the
	// destination for the data and then, piece by piece, the Handler can
	// send the data to the writer as it is being retrieved.
	InitResponse(w http.ResponseWriter) error
	// The meat of the Handler.
	Action
	// Write the output to the response (headers, body).
	WriteResponse(w http.ResponseWriter, r *http.Request) error
}

// handlerExecer manages the execution of a handler.
type handlerExecer struct {
	init HandlerInitializer
}

// serve initializes and executes the handler. The handler's methods are invoked
// in the pre-defined order, if any of the methods return an error serve will exit
// immediately and return that error, leaving the rest of the handler's methods untouched.
func (x *handlerExecer) serve(w http.ResponseWriter, r *http.Request, c context.Context) error {
	h := x.init.Init()
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

// This method is a no-op.
func (handlerbase) AuthCheck(_ *http.Request, _ context.Context) error { return nil }

// This method is a no-op.
func (handlerbase) ReadRequest(_ *http.Request, _ context.Context) error { return nil }

// This method is a no-op.
func (handlerbase) InitResponse(_ http.ResponseWriter) error { return nil }

// This method is a no-op.
func (handlerbase) WriteResponse(_ http.ResponseWriter, _ *http.Request) error { return nil }
