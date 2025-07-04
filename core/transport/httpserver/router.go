package httpserver

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(mux *http.ServeMux) *Router {
	return &Router{mux: mux}
}

func (r *Router) Post(path string, handlerFunc http.HandlerFunc) {
	r.mux.Handle("POST "+path, otelhttp.NewHandler(handlerFunc, path,
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attribute.String("resource.name", fmt.Sprintf("POST %v", path))),
		),
	))
}

func (r *Router) Get(path string, handlerFunc http.HandlerFunc) {
	r.mux.Handle("GET "+path, otelhttp.NewHandler(handlerFunc, path,
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attribute.String("resource.name", fmt.Sprintf("GET %v", path))),
		),
	))
}

func (r *Router) Put(path string, handlerFunc http.HandlerFunc) {
	r.mux.Handle("PUT "+path, otelhttp.NewHandler(handlerFunc, path,
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attribute.String("resource.name", fmt.Sprintf("PUT %v", path))),
		),
	))
}

func (r *Router) Delete(path string, handlerFunc http.HandlerFunc) {
	r.mux.Handle("DELETE "+path, otelhttp.NewHandler(handlerFunc, path,
		otelhttp.WithSpanOptions(
			trace.WithAttributes(attribute.String("resource.name", fmt.Sprintf("DELETE %v", path))),
		),
	))
}

// ServeHTTP handles HTTP requests
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
