package service

import (
	"context"
	"net/http"

	"github.com/yourorg/go-api-template/core/health"
	"github.com/yourorg/go-api-template/internal/model"
	"github.com/yourorg/go-api-template/internal/repository"
)

// HealthService provides health check functionality
type HealthServiceInterface interface {
	HealthCheck(ctx context.Context) (*model.HealthCheckResponse, error)
	Liveness(ctx context.Context) (*model.LivenessResponse, error)
	Readiness(ctx context.Context) (*model.ReadinessResponse, error)
}

type healthService struct {
	healthChecker *health.HealthService
}

// NewHealthService creates a new health service
func NewHealthService(repo *repository.Repository) HealthServiceInterface {
	healthChecker := health.NewHealthService("v1.0.0")

	// Register database checker if database is available
	if repo != nil && repo.DB != nil {
		dbChecker := health.NewPgxDatabaseChecker(repo.DB)
		healthChecker.RegisterChecker("database", dbChecker)
	}

	return &healthService{
		healthChecker: healthChecker,
	}
}

// HealthCheck performs a comprehensive health check
func (s *healthService) HealthCheck(ctx context.Context) (*model.HealthCheckResponse, error) {
	healthResult := s.healthChecker.Check(ctx)

	// Determine HTTP status based on health status
	status := http.StatusOK
	switch healthResult.Status {
	case health.StatusUnhealthy:
		status = http.StatusServiceUnavailable
	case health.StatusDegraded:
		status = http.StatusOK // Still serving traffic but degraded
	}

	return &model.HealthCheckResponse{
		Status: status,
		Data:   healthResult,
	}, nil
}

// Liveness performs a liveness check (is the application running?)
func (s *healthService) Liveness(ctx context.Context) (*model.LivenessResponse, error) {
	livenessResult := s.healthChecker.Liveness(ctx)

	return &model.LivenessResponse{
		Status: http.StatusOK, // Liveness should always be OK if we can respond
		Data:   livenessResult,
	}, nil
}

// Readiness performs a readiness check (is the application ready to serve traffic?)
func (s *healthService) Readiness(ctx context.Context) (*model.ReadinessResponse, error) {
	readinessResult := s.healthChecker.Readiness(ctx)

	// Determine HTTP status based on readiness
	status := http.StatusOK
	if readinessResult.Status == health.StatusUnhealthy {
		status = http.StatusServiceUnavailable
	}

	return &model.ReadinessResponse{
		Status: status,
		Data:   readinessResult,
	}, nil
}
