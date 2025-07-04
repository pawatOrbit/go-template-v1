package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const name = "middlewares"

var (
	tracer = otel.Tracer(name)
)

// OtelMiddleware is a middleware that automatically inject span an enrich the span with necessary metrics
func OtelMiddleware() TransportMiddleware {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewMiddleware("http.request",
			otelhttp.WithTracerProvider(otel.GetTracerProvider()),
			otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})),
		)(otelDatadogCustomAttr(next))
	}
}
func otelDatadogCustomAttr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		span := trace.SpanFromContext(r.Context())
		resourceNameKey := attribute.Key("resource.name")
		span.SetAttributes(resourceNameKey.String(r.URL.Path))
	})
}
