package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/pawatOrbit/ai-mock-data-service/go/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/exception"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/httpclient"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/logger"
	middleware_httpserver "github.com/pawatOrbit/ai-mock-data-service/go/core/transport/httpserver/middlewares"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/repository"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/service"
	"github.com/pawatOrbit/ai-mock-data-service/go/utils"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewHttpServer() (*http.Server, error) {
	cfg := config.GetConfig()
	slog.InfoContext(context.Background(), "Initializing HTTP server", "port", cfg.RestServer.Port)
	var middlewares []middleware_httpserver.TransportMiddleware
	middlewares = append(middlewares, cors.New(cors.Options{
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		AllowedMethods: cfg.CORS.AllowedMethods,
		AllowedHeaders: cfg.CORS.AllowedHeaders,
		ExposedHeaders: cfg.CORS.ExposedHeaders,
		MaxAge:         cfg.CORS.MaxAge,
	}).Handler)

	middlewareStack := middleware_httpserver.CreateStack(middlewares...)

	// Create repository
	repo, err := repository.NewRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	mockDataAppError := exception.NewMockDataServiceErrors()

	utils := utils.NewUtils()

	logger := *logger.Slog

	lmStudioClient := httpclient.NewLmStudioHttpClient(&cfg.LMStudio, logger)

	service := service.NewService(
		repo,
		cfg,
		mockDataAppError,
		utils,
		lmStudioClient,
	)

	handler := registerRoute(service)
	wrappedMiddleware := middlewareStack(handler)
	wrappedOtel := otelhttp.NewHandler(
		wrappedMiddleware,
		"",
		otelhttp.WithSpanNameFormatter(
			func(operation string, r *http.Request) string {
				return fmt.Sprintf("%s %s %s", operation, r.Method, r.URL.Path)
			},
		))

	return &http.Server{
		Addr:    ":" + cfg.RestServer.Port,
		Handler: wrappedOtel,
	}, nil
}
