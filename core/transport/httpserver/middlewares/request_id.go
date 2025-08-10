package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/yourorg/go-api-template/core/logger"
)

// RequestIDKey is the key used to store request ID in context
type RequestIDKey string

const (
	// RequestIDContextKey is the context key for request ID
	RequestIDContextKey RequestIDKey = "request_id"
	// RequestIDHeader is the default header name for request ID
	RequestIDHeader = "X-Request-ID"
	// CorrelationIDHeader is an alternative header name for request ID
	CorrelationIDHeader = "X-Correlation-ID"
	// TraceIDHeader is another alternative header name
	TraceIDHeader = "X-Trace-ID"
)

// RequestIDConfig configures the request ID middleware
type RequestIDConfig struct {
	// HeaderNames to check for existing request IDs (in order of preference)
	HeaderNames []string
	// ResponseHeader is the header name to set in the response
	ResponseHeader string
	// Generator is a function to generate new request IDs
	Generator func() string
}

// DefaultRequestIDConfig returns a default configuration
func DefaultRequestIDConfig() RequestIDConfig {
	return RequestIDConfig{
		HeaderNames:    []string{RequestIDHeader, CorrelationIDHeader, TraceIDHeader},
		ResponseHeader: RequestIDHeader,
		Generator:      generateUUIDv4,
	}
}

// RequestIDMiddleware creates a middleware that tracks request IDs
func RequestIDMiddleware(config RequestIDConfig) func(http.Handler) http.Handler {
	// Use default config if not provided
	if len(config.HeaderNames) == 0 {
		config = DefaultRequestIDConfig()
	}

	if config.Generator == nil {
		config.Generator = generateUUIDv4
	}

	if config.ResponseHeader == "" {
		config.ResponseHeader = RequestIDHeader
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestID string

			// Try to get request ID from headers (in order of preference)
			for _, headerName := range config.HeaderNames {
				if id := r.Header.Get(headerName); id != "" {
					requestID = id
					break
				}
			}

			// Generate new request ID if not found
			if requestID == "" {
				requestID = config.Generator()
			}

			// Add request ID to response header
			w.Header().Set(config.ResponseHeader, requestID)

			// Add request ID to context
			ctx := context.WithValue(r.Context(), RequestIDContextKey, requestID)

			// Add request ID to structured logging context
			ctx = logger.AddFieldsToContext(ctx, map[string]interface{}{
				"request_id": requestID,
			})

			// Continue with the enriched context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestIDFromContext extracts request ID from context
func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(RequestIDContextKey).(string)
	return requestID, ok
}

// MustGetRequestIDFromContext extracts request ID from context or returns empty string
func MustGetRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := GetRequestIDFromContext(ctx); ok {
		return requestID
	}
	return ""
}

// SetRequestIDInContext sets request ID in context
func SetRequestIDInContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDContextKey, requestID)
}

// generateUUIDv4 generates a UUID v4 as request ID
func generateUUIDv4() string {
	return uuid.New().String()
}

// generateShortID generates a shorter request ID (8 characters)
func generateShortID() string {
	id := uuid.New().String()
	// Take first 8 characters for shorter IDs
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}

// RequestIDFromHTTPHandler is a convenience wrapper for HTTP handlers
// It extracts the request ID and adds it to the handler context
func RequestIDFromHTTPHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This assumes RequestIDMiddleware has already run
		requestID := MustGetRequestIDFromContext(r.Context())

		// Log the request with ID
		logger.Slog.InfoContext(r.Context(), "Processing request",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		handler(w, r)
	}
}
