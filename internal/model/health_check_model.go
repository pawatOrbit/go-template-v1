package model

import (
	"github.com/yourorg/go-api-template/core/health"
)

// Deprecated: Use health.HealthResponse instead
type HealthReq struct {
	Name string `json:"name"`
}

// Deprecated: Use health.HealthResponse instead
type HealthResp struct {
	Status   int    `json:"status"`
	Response string `json:"response"`
}

// HealthCheckResponse wraps the health response for API
type HealthCheckResponse struct {
	Status int                   `json:"status"`
	Data   health.HealthResponse `json:"data"`
}

// LivenessResponse represents liveness check response
type LivenessResponse struct {
	Status int                   `json:"status"`
	Data   health.HealthResponse `json:"data"`
}

// ReadinessResponse represents readiness check response
type ReadinessResponse struct {
	Status int                   `json:"status"`
	Data   health.HealthResponse `json:"data"`
}
