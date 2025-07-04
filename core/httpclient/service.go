package httpclient

import (
	"log/slog"

	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/httpclient/completions"
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
