package server

import (
	"context"
	"net/http"

	"github.com/yourorg/go-api-template/core/transport/httpserver"
	middleware_httpserver "github.com/yourorg/go-api-template/core/transport/httpserver/middlewares"
	"github.com/yourorg/go-api-template/internal/model"
	"github.com/yourorg/go-api-template/internal/service"
)

func registerRoute(service service.Service) http.Handler {
	mux := http.NewServeMux()
	r := httpserver.NewRouter(mux)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware_httpserver.NotFound(w, r)
	}))

	// Health check endpoints (no authentication required)
	r.Get("/health", httpserver.NewTransport(
		&struct{}{},
		httpserver.NewEndpoint(func(ctx context.Context, in *struct{}) (*model.HealthCheckResponse, error) {
			return service.HealthService.HealthCheck(ctx)
		}),
	))

	r.Get("/health/liveness", httpserver.NewTransport(
		&struct{}{},
		httpserver.NewEndpoint(func(ctx context.Context, in *struct{}) (*model.LivenessResponse, error) {
			return service.HealthService.Liveness(ctx)
		}),
	))

	r.Get("/health/readiness", httpserver.NewTransport(
		&struct{}{},
		httpserver.NewEndpoint(func(ctx context.Context, in *struct{}) (*model.ReadinessResponse, error) {
			return service.HealthService.Readiness(ctx)
		}),
	))

	// Example API endpoints - replace with your actual endpoints
	r.Get("/api/v1/examples/{id}", httpserver.NewTransport(
		&model.ExampleRequest{},
		httpserver.NewEndpoint(service.ExampleService.GetExample),
	))

	r.Post("/api/v1/examples", httpserver.NewTransport(
		&model.CreateExampleRequest{},
		httpserver.NewEndpoint(service.ExampleService.CreateExample),
	))

	// Legacy health check endpoint (deprecated)
	r.Post("/health-check",
		httpserver.NewTransport(
			&model.HealthReq{},
			httpserver.NewEndpoint(func(ctx context.Context, in *model.HealthReq) (*model.HealthResp, error) {
				return &model.HealthResp{
					Status:   1000,
					Response: "Hello, " + in.Name,
				}, nil
			})))
	return mux
}
