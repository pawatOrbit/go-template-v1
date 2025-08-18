package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/yourorg/go-api-template/core/exception"
	"github.com/yourorg/go-api-template/core/logger"
	"github.com/yourorg/go-api-template/core/transport"
	middleware "github.com/yourorg/go-api-template/core/transport/httpserver/middlewares"
)

type errorResp struct {
	Status       int               `json:"status"`
	Message      string            `json:"message"`
	DebugMessage string            `json:"debug_message,omitempty"`
	Fields       []string          `json:"fields,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
}

func NewTransport[T, R any](req T, endpoint func() Endpoint[T, R], middlewares ...transport.EndpointMiddleware[T, R]) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		newReq := deepCopy(req)

		var (
			ctx            = r.Context()
			httpStatusCode = http.StatusOK
			method         = r.Method
			path           = r.URL.Path
			header         = r.Header
			requestBody    []byte
			resp           R
			serviceError   error
			elapsedTime    time.Duration
		)

		requestBody, err := readRequestBody(r)
		if err != nil {
			fmt.Println("Error reading request body")
			HandleInternalServerError(w, http.StatusBadRequest)
			return
		}

		if len(requestBody) == 0 {
			requestBody = []byte("{}")
		}

		err = json.Unmarshal(requestBody, &newReq)
		if err != nil {
			fmt.Println("Error unmarshalling request body")
			HandleInternalServerError(w, http.StatusBadRequest)
			return
		}

		startTime := time.Now()
		resp, serviceError = endpoint()()(r.Context(), newReq)
		elapsedTime = time.Since(startTime)

		if serviceError != nil {
			// Check if error is an ExceptionError to get proper status code
			if exErr, ok := serviceError.(*exception.ExceptionError); ok {
				httpStatusCode = exErr.HttpStatusCode
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpStatusCode)
				// Build error response with all available fields
				errorResponse := errorResp{
					Status:  exErr.APIStatusCode,
					Message: exErr.GlobalMessage,
					Fields:  exErr.ErrFields,
					Data:    exErr.ErrWithDatas,
				}
				json.NewEncoder(w).Encode(errorResponse)
			} else {
				httpStatusCode = http.StatusInternalServerError
				HandleInternalServerError(w, httpStatusCode)
			}
			logRequestAndResponse(ctx, startTime, elapsedTime, method, path, header, requestBody, []byte(fmt.Sprintf("%v", resp)), serviceError, httpStatusCode)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpStatusCode)
			json.NewEncoder(w).Encode(resp)

			logRequestAndResponse(ctx, startTime, elapsedTime, method, path, header, requestBody, []byte(fmt.Sprintf("%v", resp)), serviceError, httpStatusCode)
			return
		}
	}
}

func deepCopy[T any](src T) T {
	return reflect.New(reflect.TypeOf(src).Elem()).Interface().(T)
}

func readRequestBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func logRequestAndResponse(ctx context.Context, startTime time.Time, elapsedTime time.Duration, method string, path string, header http.Header, requestBody, responseBody []byte, serviceError error, httpStatusCode int) {
	middleware.LoggingNetHttp(ctx, *logger.Slog, startTime, elapsedTime, method, path, header, requestBody, responseBody, serviceError, httpStatusCode)
}
