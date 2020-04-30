package hxio

import (
	"context"
	"net/http"

	"github.com/frk/route"
)

type ErrorHandler interface {
	HandleError(w http.ResponseWriter, r *http.Request, err error)
}

type RouteOptions struct {
	HandlerInitializerFactory HandlerInitializerFactory
	ErrorHandler              ErrorHandler
	PathPrefix                string
}

type RouteList []struct {
	Path    string
	Method  string
	Handler interface{}
}

func InitRouter(r *route.Router, routes RouteList, opts RouteOptions) {
	if opts.HandlerInitializerFactory == nil {
		opts.HandlerInitializerFactory = handlerInitializerFactory{}
	}
	if opts.ErrorHandler == nil {
		opts.ErrorHandler = errorHandler{}
	}

	for _, rt := range routes {
		path := opts.PathPrefix + rt.Path
		method := rt.Method

		handler := new(routeHandler)
		handler.init = opts.HandlerInitializerFactory.NewHandlerInitializer(rt.Handler, path, method)
		handler.eh = opts.ErrorHandler

		r.Handle(method, path, handler)
	}
}

type routeHandler struct {
	handlerExecer
	eh ErrorHandler
}

func (h *routeHandler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := h.serve(w, r, ctx); err != nil {
		h.eh.HandleError(w, r, err)
	}
}

func InitServeMux(m *http.ServeMux, routes RouteList, opts RouteOptions) {
	if opts.HandlerInitializerFactory == nil {
		opts.HandlerInitializerFactory = handlerInitializerFactory{}
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
		handler.init = opts.HandlerInitializerFactory.NewHandlerInitializer(rt.Handler, path, method)
		handler.eh = opts.ErrorHandler

		mux.Handle(path, handler)
	}
}

type httpHandler struct {
	handlerExecer
	eh ErrorHandler
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.serve(w, r, r.Context()); err != nil {
		h.eh.HandleError(w, r, err)
	}
}

type errorHandler struct{}

func (errorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

type handlerInitializerFactory struct{}

func (handlerInitializerFactory) NewHandlerInitializer(v interface{}, path, method string) HandlerInitializer {
	return v.(HandlerInitializer)
}
