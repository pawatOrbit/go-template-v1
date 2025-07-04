package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/pawatOrbit/ai-mock-data-service/go/core/logger"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/transport"
	middleware "github.com/pawatOrbit/ai-mock-data-service/go/core/transport/httpserver/middlewares"
)

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
			httpStatusCode = http.StatusInternalServerError
			HandleInternalServerError(w, httpStatusCode)
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
