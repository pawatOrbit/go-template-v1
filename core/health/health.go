package health

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/go-api-template/core/logger"
)

// Status represents the health status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// ComponentHealth represents the health of a single component
type ComponentHealth struct {
	Name      string            `json:"name"`
	Status    Status            `json:"status"`
	Message   string            `json:"message,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Duration  time.Duration     `json:"duration"`
}

// HealthResponse represents the overall health response
type HealthResponse struct {
	Status     Status                     `json:"status"`
	Timestamp  time.Time                  `json:"timestamp"`
	Version    string                     `json:"version,omitempty"`
	Components map[string]ComponentHealth `json:"components"`
	System     SystemInfo                 `json:"system"`
}

// SystemInfo represents system information
type SystemInfo struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	NumCPU    int    `json:"num_cpu"`
	NumGoroutines int    `json:"num_goroutines"`
	MemStats  MemoryStats `json:"memory"`
}

// MemoryStats represents memory statistics
type MemoryStats struct {
	Alloc        uint64 `json:"alloc_mb"`
	TotalAlloc   uint64 `json:"total_alloc_mb"`
	Sys          uint64 `json:"sys_mb"`
	NumGC        uint32 `json:"num_gc"`
}

// Checker interface for health checks
type Checker interface {
	Check(ctx context.Context) ComponentHealth
}

// HealthService manages health checks
type HealthService struct {
	checkers map[string]Checker
	version  string
}

// NewHealthService creates a new health service
func NewHealthService(version string) *HealthService {
	return &HealthService{
		checkers: make(map[string]Checker),
		version:  version,
	}
}

// RegisterChecker registers a health checker
func (hs *HealthService) RegisterChecker(name string, checker Checker) {
	hs.checkers[name] = checker
}

// Check performs all health checks
func (hs *HealthService) Check(ctx context.Context) HealthResponse {
	start := time.Now()
	components := make(map[string]ComponentHealth)
	overallStatus := StatusHealthy

	// Run all registered health checks
	for name, checker := range hs.checkers {
		componentHealth := checker.Check(ctx)
		components[name] = componentHealth

		// Determine overall status
		if componentHealth.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
		} else if componentHealth.Status == StatusDegraded && overallStatus != StatusUnhealthy {
			overallStatus = StatusDegraded
		}
	}

	return HealthResponse{
		Status:     overallStatus,
		Timestamp:  start,
		Version:    hs.version,
		Components: components,
		System:     getSystemInfo(),
	}
}

// Liveness performs a basic liveness check (application is running)
func (hs *HealthService) Liveness(ctx context.Context) HealthResponse {
	return HealthResponse{
		Status:    StatusHealthy,
		Timestamp: time.Now(),
		Version:   hs.version,
		System:    getSystemInfo(),
	}
}

// Readiness performs readiness check (application is ready to serve traffic)
func (hs *HealthService) Readiness(ctx context.Context) HealthResponse {
	// For readiness, we typically check critical dependencies
	criticalComponents := make(map[string]ComponentHealth)
	overallStatus := StatusHealthy

	// Check only critical components for readiness
	for name, checker := range hs.checkers {
		// Only check database for readiness (Redis is not critical for serving traffic)
		if name == "database" {
			componentHealth := checker.Check(ctx)
			criticalComponents[name] = componentHealth

			if componentHealth.Status == StatusUnhealthy {
				overallStatus = StatusUnhealthy
			}
		}
	}

	return HealthResponse{
		Status:     overallStatus,
		Timestamp:  time.Now(),
		Version:    hs.version,
		Components: criticalComponents,
		System:     getSystemInfo(),
	}
}

// getSystemInfo returns system information
func getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		Version:       "1.0.0", // This could be injected during build
		GoVersion:     runtime.Version(),
		NumCPU:        runtime.NumCPU(),
		NumGoroutines: runtime.NumGoroutine(),
		MemStats: MemoryStats{
			Alloc:      bToMb(m.Alloc),
			TotalAlloc: bToMb(m.TotalAlloc),
			Sys:        bToMb(m.Sys),
			NumGC:      m.NumGC,
		},
	}
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// DatabaseChecker checks database connectivity for standard sql.DB
type DatabaseChecker struct {
	db *sql.DB
}

// NewDatabaseChecker creates a new database checker
func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{db: db}
}

// Check implements the Checker interface for database
func (dc *DatabaseChecker) Check(ctx context.Context) ComponentHealth {
	start := time.Now()
	
	if dc.db == nil {
		return ComponentHealth{
			Name:      "database",
			Status:    StatusUnhealthy,
			Message:   "Database connection is nil",
			Timestamp: start,
			Duration:  time.Since(start),
		}
	}

	// Create a context with timeout for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := dc.db.PingContext(pingCtx)
	duration := time.Since(start)

	if err != nil {
		logger.Slog.Error("Database health check failed", "error", err.Error())
		return ComponentHealth{
			Name:      "database",
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Database ping failed: %v", err),
			Timestamp: start,
			Duration:  duration,
		}
	}

	// Check connection stats
	stats := dc.db.Stats()
	details := map[string]string{
		"open_connections":     fmt.Sprintf("%d", stats.OpenConnections),
		"max_open_connections": fmt.Sprintf("%d", stats.MaxOpenConnections),
		"idle_connections":     fmt.Sprintf("%d", stats.Idle),
	}

	status := StatusHealthy
	message := "Database is healthy"

	// Check if we're running low on connections
	if stats.MaxOpenConnections > 0 && float64(stats.OpenConnections)/float64(stats.MaxOpenConnections) > 0.8 {
		status = StatusDegraded
		message = "Database connection pool is running low"
	}

	return ComponentHealth{
		Name:      "database",
		Status:    status,
		Message:   message,
		Details:   details,
		Timestamp: start,
		Duration:  duration,
	}
}

// PgxDatabaseChecker checks pgx database connectivity
type PgxDatabaseChecker struct {
	pool *pgxpool.Pool
}

// NewPgxDatabaseChecker creates a new pgx database checker
func NewPgxDatabaseChecker(pool *pgxpool.Pool) *PgxDatabaseChecker {
	return &PgxDatabaseChecker{pool: pool}
}

// Check implements the Checker interface for pgx database
func (pdc *PgxDatabaseChecker) Check(ctx context.Context) ComponentHealth {
	start := time.Now()
	
	if pdc.pool == nil {
		return ComponentHealth{
			Name:      "database",
			Status:    StatusUnhealthy,
			Message:   "Database pool is nil",
			Timestamp: start,
			Duration:  time.Since(start),
		}
	}

	// Create a context with timeout for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := pdc.pool.Ping(pingCtx)
	duration := time.Since(start)

	if err != nil {
		logger.Slog.Error("Database health check failed", "error", err.Error())
		return ComponentHealth{
			Name:      "database",
			Status:    StatusUnhealthy,
			Message:   fmt.Sprintf("Database ping failed: %v", err),
			Timestamp: start,
			Duration:  duration,
		}
	}

	// Check connection stats
	stats := pdc.pool.Stat()
	details := map[string]string{
		"total_connections":       fmt.Sprintf("%d", stats.TotalConns()),
		"idle_connections":        fmt.Sprintf("%d", stats.IdleConns()),
		"acquired_connections":    fmt.Sprintf("%d", stats.AcquiredConns()),
		"constructing_connections": fmt.Sprintf("%d", stats.ConstructingConns()),
		"max_connections":         fmt.Sprintf("%d", stats.MaxConns()),
	}

	status := StatusHealthy
	message := "Database is healthy"

	// Check if we're running low on connections (80% threshold)
	if stats.MaxConns() > 0 && float64(stats.AcquiredConns())/float64(stats.MaxConns()) > 0.8 {
		status = StatusDegraded
		message = "Database connection pool is running low"
	}

	return ComponentHealth{
		Name:      "database",
		Status:    status,
		Message:   message,
		Details:   details,
		Timestamp: start,
		Duration:  duration,
	}
}