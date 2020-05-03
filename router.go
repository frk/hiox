package hiox

import (
	"context"
	"net/http"

	"github.com/frk/route"
)

// HandlerInitializerAdapter converts a custom, user-defined, handler initializer
// into a value that implements the hiox.HandlerInitializer interface.
type HandlerInitializerAdapter interface {
	AdaptHandlerInitializer(hi interface{}, path, method string) HandlerInitializer
}

// ErrorHandler interface is used to handle the errors returned from Handlers.
// When one of the Handler's methods returns an error, that error will be passed
// in to the ErrorHandler which then has a chance to format the error response
// as it sees fit.
//
// If no ErrorHandler is provided, then by default the http.Error function is
// used to write the response using the err.Error() as the response text and
// the http.StatusBadRequest as the status code.
type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

// RouteOptions is a set of options that, if set, will be applied
// to each route being registered.
type RouteOptions struct {
	// The HandlerInitializerAdapter to be used to convert the custom
	// route Handler value into an hiox.Handler.
	HandlerInitializerAdapter HandlerInitializerAdapter
	// The ErrorHandler to be used to handle errors returned
	// from the route Handlers.
	ErrorHandler ErrorHandler
	// The prefix to be applied to the routes' paths.
	PathPrefix string
}

// RouteList is a list of settings used to register HandlerInitializers for the specified paths.
type RouteList []struct {
	// The path pattern for which to register the HandlerInitializer.
	Path string
	// The HTTP method for which to register the HandlerInitializer.
	Method string
	// The HandlerInitializer to be registered.
	HandlerInitializer interface{}
}

// InitRouter takes the HandlerInitializers in the provided RouteList
// and registers them as route.Handlers with the given *route.Router.
func InitRouter(r *route.Router, routes RouteList, opts RouteOptions) {
	if opts.HandlerInitializerAdapter == nil {
		opts.HandlerInitializerAdapter = handlerInitializerAdapter{}
	}
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = errorHandler{}
	}

	for _, rt := range routes {
		path := opts.PathPrefix + rt.Path
		method := rt.Method

		handler := new(routeHandler)
		handler.init = opts.HandlerInitializerAdapter.AdaptHandlerInitializer(rt.HandlerInitializer, path, method)
		handler.eh = opts.ErrorHandler

		r.Handle(method, path, handler)
	}
}

// routeHandler is a wrapper around handlerExecer that implements the route.Handler interface.
type routeHandler struct {
	handlerExecer
	eh ErrorHandler
}

func (h *routeHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.serve(w, r, ctx); err != nil {
		h.eh.HandleError(w, r, err)
	}
}

// InitServeMux takes the HandlerInitializers in the provided RouteList
// and registers them as http.Handlers with the given *http.ServeMux.
//
// NOTE(mkopriva): InitServeMux can be called only once per *http.ServeMux.
func InitServeMux(m *http.ServeMux, routes RouteList, opts RouteOptions) {
	if opts.HandlerInitializerAdapter == nil {
		opts.HandlerInitializerAdapter = handlerInitializerAdapter{}
	}
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = errorHandler{}
	}

	muxmap := make(map[string]*http.ServeMux)
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if mux, ok := muxmap[r.Method]; ok {
			mux.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})

	for _, rt := range routes {
		path := opts.PathPrefix + rt.Path
		method := rt.Method

		mux, ok := muxmap[method]
		if !ok {
			mux = http.NewServeMux()
			muxmap[method] = mux
		}

		handler := new(httpHandler)
		handler.init = opts.HandlerInitializerAdapter.AdaptHandlerInitializer(rt.HandlerInitializer, path, method)
		handler.eh = opts.ErrorHandler

		mux.Handle(path, handler)
	}
}

// httpHandler is a wrapper around handlerExecer that implements the http.Handler interface.
type httpHandler struct {
	handlerExecer
	eh ErrorHandler
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.serve(w, r, r.Context()); err != nil {
		h.eh.HandleError(w, r, err)
	}
}

// Default ErrorHandler implementation.
type errorHandler struct{}

func (errorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

// Default HandlerInitializerAdapter implementation.
type handlerInitializerAdapter struct{}

func (handlerInitializerAdapter) AdaptHandlerInitializer(v interface{}, path, method string) HandlerInitializer {
	return v.(HandlerInitializer)
}
