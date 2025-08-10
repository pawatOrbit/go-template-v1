package httpclient

import (
	"log/slog"

	core_config "github.com/yourorg/go-api-template/core/config"
	"github.com/yourorg/go-api-template/core/httpclient/completions"
)

type LmStudioServiceClient struct {
	GetCompletionsService completions.CompletionsServiceClient
}

func NewLmStudioHttpClient(cfg *core_config.LMStudioConfig, logger slog.Logger) *LmStudioServiceClient {
	return &LmStudioServiceClient{
		GetCompletionsService: completions.NewCompletionsServiceClient(cfg, logger),
	}
}

func NewMockLmStudioClient(cfg *core_config.LMStudioConfig) *LmStudioServiceClient {
	return &LmStudioServiceClient{}
}
