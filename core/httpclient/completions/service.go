package completions

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/httpclient/common"
)

type CompletionsServiceClient interface {
	GetCompletionsService(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

type completionsServiceClient struct {
	cfg        *core_config.LMStudioConfig
	httpClient *http.Client
	logger     slog.Logger
}

func NewCompletionsServiceClient(cfg *core_config.LMStudioConfig, logger slog.Logger) CompletionsServiceClient {
	httpClient := http.Client{
		Timeout: 10 * time.Minute,
	}
	return &completionsServiceClient{
		cfg:        cfg,
		httpClient: &httpClient,
		logger:     logger,
	}
}

func (s *completionsServiceClient) GetCompletionsService(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	path := GET_COMPLETIONS_URL
	slogger := s.logger.With("method", "GetCompletionsService")
	return common.Do[CompletionRequest, CompletionResponse, *CompletionError](ctx, s.cfg, s.httpClient, path, req, slogger)
}
