package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/logger"
	"go.opentelemetry.io/otel"
)

const name = "aoa.common"

var tracer = otel.Tracer(name)

func Do[T any, R any, E error](ctx context.Context, cfg *core_config.LMStudioConfig, httpClient *http.Client, path string, req interface{}, slogger *slog.Logger) (R, error) {
	ctx, span := tracer.Start(ctx, path)
	defer span.End()

	typedResp := new(R)

	typedReq, ok := req.(T)
	if !ok {
		return *typedResp, errors.New("req is not a valid type")
	}

	starTime := time.Now()
	if cfg == nil {
		return *typedResp, errors.New("AOA config is nil")
	}

	payload, err := json.Marshal(typedReq)
	if err != nil {
		return *typedResp, err
	}

	fullPath, err := BuildURL(cfg, path) // common lms for build full url
	if err != nil {
		return *typedResp, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fullPath, bytes.NewReader(payload))
	if err != nil {
		return *typedResp, err
	}

	token := Basic
	httpResp, err := DefaultDo(ctx, cfg, r, httpClient, ApplicationJson, token, nil)
	if err != nil {
		return *typedResp, err
	}
	defer httpResp.Body.Close()

	responseData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return *typedResp, err
	}

	var commonErrorResponse error
	respErr := new(E)
	logLevel := logger.Info

	switch httpResp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		err = json.Unmarshal(responseData, typedResp)
		if err != nil {
			fmt.Println(err.Error())
			return *typedResp, err
		}
	case http.StatusBadRequest, http.StatusInternalServerError:
		err = json.Unmarshal(responseData, respErr)
		if err != nil {
			return *typedResp, err
		}
		logLevel = logger.Error
		commonErrorResponse = *respErr
	default:
		commonErrorResponse = TransportError{
			Code:        httpResp.StatusCode,
			Description: fmt.Sprintf("got %d response from %s is %s", httpResp.StatusCode, fullPath, responseData),
		}
		logLevel = logger.Error
	}

	logger.CanonicalLogger(
		ctx,
		*slogger,
		logLevel,
		payload,
		responseData,
		commonErrorResponse,
		logger.CanonicalLog{
			Transport: "http",
			Traffic:   "external",
			Method:    http.MethodPost,
			Status:    httpResp.StatusCode,
			Path:      path,
			Duration:  time.Since(starTime),
		},
		[]any{},
	)

	return *typedResp, commonErrorResponse
}
