package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
	"github.com/yourorg/go-api-template/config"
	"github.com/yourorg/go-api-template/core/cache"
	"github.com/yourorg/go-api-template/core/exception"
	"github.com/yourorg/go-api-template/core/httpclient"
	"github.com/yourorg/go-api-template/core/logger"
	"github.com/yourorg/go-api-template/core/ratelimit"
	middleware_httpserver "github.com/yourorg/go-api-template/core/transport/httpserver/middlewares"
	"github.com/yourorg/go-api-template/internal/repository"
	"github.com/yourorg/go-api-template/internal/service"
	"github.com/yourorg/go-api-template/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewHttpServer() (*http.Server, error) {
	cfg := config.GetConfig()
	slog.InfoContext(context.Background(), "Initializing HTTP server", "port", cfg.RestServer.Port)
	var middlewares []middleware_httpserver.TransportMiddleware

	// CORS middleware
	middlewares = append(middlewares, cors.New(cors.Options{
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		AllowedMethods: cfg.CORS.AllowedMethods,
		AllowedHeaders: cfg.CORS.AllowedHeaders,
		ExposedHeaders: cfg.CORS.ExposedHeaders,
		MaxAge:         cfg.CORS.MaxAge,
	}).Handler)

	// Rate limiting middleware
	if cfg.RateLimit.Enabled {
		// Initialize Redis cache service for rate limiting
		err := cache.InitRedisService(cfg.Redis)
		if err != nil {
			slog.WarnContext(context.Background(), "Failed to initialize Redis for rate limiting, using memory limiter", "error", err.Error())
		}

		// Create rate limiter based on available cache service
		var limiter ratelimit.Limiter
		if cacheService := cache.GetRedisService(); cacheService != nil {
			limiter = ratelimit.NewRedisLimiter(cacheService, createRateLimitConfig(cfg))
			slog.InfoContext(context.Background(), "Using Redis-based rate limiter")
		} else {
			limiter = ratelimit.NewMemoryLimiter(createRateLimitConfig(cfg))
			slog.InfoContext(context.Background(), "Using memory-based rate limiter")
		}

		middlewares = append(middlewares, ratelimit.Middleware(limiter, createRateLimitConfig(cfg)))
		slog.InfoContext(context.Background(), "Rate limiting enabled",
			"requests", cfg.RateLimit.Requests,
			"window", cfg.RateLimit.Window)
	}

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

// createRateLimitConfig converts config values to ratelimit.Config
func createRateLimitConfig(cfg *config.Config) ratelimit.Config {
	window, err := time.ParseDuration(cfg.RateLimit.Window)
	if err != nil {
		slog.WarnContext(context.Background(), "Invalid rate limit window duration, using default", "window", cfg.RateLimit.Window, "error", err.Error())
		window = time.Hour
	}

	config := ratelimit.Config{
		Requests:       cfg.RateLimit.Requests,
		Window:         window,
		KeyBuilder:     ratelimit.DefaultKeyBuilder,
		SkipPaths:      cfg.RateLimit.SkipPaths,
		IncludeHeaders: cfg.RateLimit.IncludeHeaders,
		Message:        cfg.RateLimit.Message,
		StatusCode:     cfg.RateLimit.StatusCode,
	}

	// Set defaults if not configured
	if config.Requests == 0 {
		config.Requests = 100
	}
	if config.Message == "" {
		config.Message = "Rate limit exceeded"
	}
	if config.StatusCode == 0 {
		config.StatusCode = http.StatusTooManyRequests
	}
	if len(config.SkipPaths) == 0 {
		config.SkipPaths = []string{"/health", "/health/*", "/metrics"}
	}

	return config
}
