package common

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
)

type ContentType int

const (
	ApplicationJson ContentType = 1 << iota
	XWWWFormUrlencoded
	MultiPartFormData
)

type TokenType int

const (
	Bearer TokenType = 1 << iota
	Basic
)

func GetCommonHeaderFromRequest(r *http.Request) []any {
	headers := []any{
		slog.String("Authorization", r.Header.Get("Authorization")),
		slog.String("RequestID", r.Header.Get("RequestID")),
		slog.String("Content-Type", r.Header.Get("Content-Type")),
	}
	return headers
}

func BuildURL(cfg *core_config.LMStudioConfig, path string) (string, error) {
	if cfg == nil {
		return "", errors.New("config is nil")
	}

	if cfg.Protocol == "" || cfg.BaseUrl == "" || path == "" {
		return "", errors.New("aoa url,protocol, and/or path are empty")
	}

	return fmt.Sprintf("%s://%s/%s", cfg.Protocol, cfg.BaseUrl, path), nil
}

func DefaultDo(ctx context.Context, cfg *core_config.LMStudioConfig, r *http.Request, c *http.Client, contentType ContentType, tokenType TokenType, customContentType *string) (*http.Response, error) {
	switch contentType {
	case ApplicationJson:
		r.Header.Set("Content-Type", "application/json")
	case XWWWFormUrlencoded:
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	case MultiPartFormData:
		r.Header.Set("Content-Type", *customContentType)
	}
	// Request ID
	reqId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	r.Header.Set("RequestID", reqId.String())
	resp, err := c.Do(r)
	if err != nil {
		print(err.Error())
		return nil, err
	}
	return resp, nil
}
