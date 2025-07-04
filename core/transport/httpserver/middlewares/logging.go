package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/pawatOrbit/ai-mock-data-service/go/core/logger"
)

// LoggingNetHttp is a middleware that logs the incoming request and the outgoing response.
// It logs the request method, path, ip, duration, accept-language, x-request-id, x-username, x-user-id, x-permissions, and the response status code.
// It also logs the request body, response body, and the error if any.
// The log level is set to error if the status code is greater than or equal to http.StatusBadRequest, otherwise it is set to info.
func LoggingNetHttp(ctx context.Context, l slog.Logger, startTime time.Time, elapse time.Duration, method string, path string, headers http.Header, requestBody []byte, responseBody []byte, err error, statusCode int) {
	ctx, span := tracer.Start(ctx, "LoggingNetHttp")
	defer span.End()

	var fields []any
	// append md log
	fields = append(fields,
		slog.String("logger_name", "canonical"),
		slog.Group("httpserver_md",
			slog.String("type", "httpserver"),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", fmt.Sprint(headers["X-Forwarded-For"])),
			slog.String("duration", elapse.String()),
		),
	)

	var level logger.Level
	if statusCode >= http.StatusBadRequest {
		level = logger.Error
		//span.SetStatus(codes.Error, err.Error())
	} else {
		level = logger.Info
		//span.SetStatus(codes.Ok, "OK")
	}

	logger.CanonicalLogger(
		ctx,
		l,
		level,
		requestBody,
		responseBody,
		err,
		logger.CanonicalLog{
			Transport: "http",
			Traffic:   "internal",
			Method:    method,
			Status:    statusCode,
			Path:      path,
			Duration:  elapse,
		},
		fields,
	)
}

func convertHeaderAttrToString(key string, headers map[string][]string) string {
	if header, ok := headers[key]; ok {
		return header[0]
	}
	return ""
}
