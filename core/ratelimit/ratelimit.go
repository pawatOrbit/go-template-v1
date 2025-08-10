package ratelimit

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yourorg/go-api-template/core/cache"
	"github.com/yourorg/go-api-template/core/logger"
)

// Config holds rate limiting configuration
type Config struct {
	// Requests per window
	Requests int `mapstructure:"requests"`
	// Window duration
	Window time.Duration `mapstructure:"window"`
	// Skip rate limiting for these paths
	SkipPaths []string `mapstructure:"skipPaths"`
	// Custom key builder function
	KeyBuilder func(*http.Request) string
	// Custom skip function
	SkipFunc func(*http.Request) bool
	// Headers to include in response
	IncludeHeaders bool `mapstructure:"includeHeaders"`
	// Message to return when rate limited
	Message string `mapstructure:"message"`
	// Status code to return when rate limited
	StatusCode int `mapstructure:"statusCode"`
}

// DefaultConfig returns default rate limiting configuration
func DefaultConfig() Config {
	return Config{
		Requests:       100,
		Window:         time.Hour,
		KeyBuilder:     DefaultKeyBuilder,
		SkipPaths:      []string{"/health", "/health/*", "/metrics"},
		IncludeHeaders: true,
		Message:        "Rate limit exceeded",
		StatusCode:     http.StatusTooManyRequests,
	}
}

// Limiter interface for rate limiting
type Limiter interface {
	Allow(ctx context.Context, key string) (bool, *Result, error)
	Reset(ctx context.Context, key string) error
}

// Result contains rate limit information
type Result struct {
	Allowed      bool
	Limit        int
	Remaining    int
	ResetTime    time.Time
	RetryAfter   time.Duration
}

// redisLimiter implements sliding window rate limiting using Redis
type redisLimiter struct {
	cacheService cache.CacheService
	config       Config
}

// memoryLimiter implements in-memory rate limiting with sliding window
type memoryLimiter struct {
	config Config
	store  map[string]*windowCounter
	mutex  sync.RWMutex
}

type windowCounter struct {
	requests []time.Time
	mutex    sync.RWMutex
}

// NewRedisLimiter creates a Redis-based rate limiter
func NewRedisLimiter(cacheService cache.CacheService, config Config) Limiter {
	return &redisLimiter{
		cacheService: cacheService,
		config:       config,
	}
}

// NewMemoryLimiter creates an in-memory rate limiter
func NewMemoryLimiter(config Config) Limiter {
	return &memoryLimiter{
		config: config,
		store:  make(map[string]*windowCounter),
	}
}

// Allow checks if request should be allowed (Redis implementation)
func (r *redisLimiter) Allow(ctx context.Context, key string) (bool, *Result, error) {
	now := time.Now()
	windowStart := now.Add(-r.config.Window)
	
	// Use Redis sorted set for sliding window
	pipe := r.cacheService.GetClient().Pipeline()
	
	// Remove old entries
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixNano()))
	
	// Count current requests
	countCmd := pipe.ZCard(ctx, key)
	
	// Add current request
	pipe.ZAdd(ctx, key, struct {
		Score  float64
		Member interface{}
	}{Score: float64(now.UnixNano()), Member: now.UnixNano()})
	
	// Set expiration
	pipe.Expire(ctx, key, r.config.Window+time.Minute)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, nil, fmt.Errorf("redis pipeline error: %w", err)
	}
	
	currentCount := int(countCmd.Val())
	allowed := currentCount < r.config.Requests
	
	result := &Result{
		Allowed:    allowed,
		Limit:      r.config.Requests,
		Remaining:  max(0, r.config.Requests-currentCount-1),
		ResetTime:  now.Add(r.config.Window),
		RetryAfter: r.config.Window,
	}
	
	if !allowed {
		result.RetryAfter = r.config.Window
	}
	
	return allowed, result, nil
}

// Allow checks if request should be allowed (Memory implementation)
func (m *memoryLimiter) Allow(ctx context.Context, key string) (bool, *Result, error) {
	now := time.Now()
	windowStart := now.Add(-m.config.Window)
	
	m.mutex.Lock()
	counter, exists := m.store[key]
	if !exists {
		counter = &windowCounter{requests: make([]time.Time, 0)}
		m.store[key] = counter
	}
	m.mutex.Unlock()
	
	counter.mutex.Lock()
	defer counter.mutex.Unlock()
	
	// Remove old requests outside the window
	validRequests := make([]time.Time, 0, len(counter.requests))
	for _, req := range counter.requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}
	counter.requests = validRequests
	
	// Check if we can allow this request
	allowed := len(counter.requests) < m.config.Requests
	
	if allowed {
		counter.requests = append(counter.requests, now)
	}
	
	result := &Result{
		Allowed:   allowed,
		Limit:     m.config.Requests,
		Remaining: max(0, m.config.Requests-len(counter.requests)),
		ResetTime: now.Add(m.config.Window),
	}
	
	if !allowed && len(counter.requests) > 0 {
		// Calculate retry after based on oldest request
		oldestRequest := counter.requests[0]
		result.RetryAfter = oldestRequest.Add(m.config.Window).Sub(now)
	}
	
	return allowed, result, nil
}

// Reset removes rate limit data for a key (Redis implementation)
func (r *redisLimiter) Reset(ctx context.Context, key string) error {
	return r.cacheService.Delete(ctx, key)
}

// Reset removes rate limit data for a key (Memory implementation)
func (m *memoryLimiter) Reset(ctx context.Context, key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.store, key)
	return nil
}

// Middleware creates rate limiting middleware
func Middleware(limiter Limiter, config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if we should skip rate limiting
			if shouldSkip(r, config) {
				next.ServeHTTP(w, r)
				return
			}
			
			// Build rate limit key
			key := config.KeyBuilder(r)
			
			// Check rate limit
			ctx := r.Context()
			allowed, result, err := limiter.Allow(ctx, key)
			if err != nil {
				if logger.Slog != nil {
					logger.Slog.ErrorContext(ctx, "Rate limiting error", "key", key, "error", err.Error())
				}
				// On error, allow the request to proceed
				next.ServeHTTP(w, r)
				return
			}
			
			// Add rate limit headers
			if config.IncludeHeaders {
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
				w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
				w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime.Unix(), 10))
				
				if !allowed {
					w.Header().Set("Retry-After", strconv.FormatInt(int64(result.RetryAfter.Seconds()), 10))
				}
			}
			
			if !allowed {
				// Rate limit exceeded
				if logger.Slog != nil {
					logger.Slog.WarnContext(ctx, "Rate limit exceeded", 
						"key", key, 
						"limit", result.Limit,
						"path", r.URL.Path,
						"ip", GetClientIP(r))
				}
				
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(config.StatusCode)
				fmt.Fprintf(w, `{"error": "%s", "retry_after": %d}`, 
					config.Message, int64(result.RetryAfter.Seconds()))
				return
			}
			
			// Request allowed, log if configured
			if logger.Slog != nil {
				logger.Slog.DebugContext(ctx, "Rate limit check passed",
					"key", key,
					"remaining", result.Remaining,
					"limit", result.Limit)
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// shouldSkip determines if rate limiting should be skipped
func shouldSkip(r *http.Request, config Config) bool {
	// Check custom skip function
	if config.SkipFunc != nil && config.SkipFunc(r) {
		return true
	}
	
	// Check skip paths
	for _, path := range config.SkipPaths {
		if path == r.URL.Path || (strings.HasSuffix(path, "/*") && 
			strings.HasPrefix(r.URL.Path, strings.TrimSuffix(path, "/*"))) {
			return true
		}
	}
	
	return false
}

// DefaultKeyBuilder builds rate limit key from client IP
func DefaultKeyBuilder(r *http.Request) string {
	ip := GetClientIP(r)
	return fmt.Sprintf("rate_limit:ip:%s", ip)
}

// UserKeyBuilder builds rate limit key from user ID (requires auth)
func UserKeyBuilder(r *http.Request) string {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		// Fallback to IP if no user ID
		return DefaultKeyBuilder(r)
	}
	return fmt.Sprintf("rate_limit:user:%s", userID)
}

// PathBasedKeyBuilder builds rate limit key including path
func PathBasedKeyBuilder(r *http.Request) string {
	ip := GetClientIP(r)
	path := r.URL.Path
	return fmt.Sprintf("rate_limit:ip:%s:path:%s", ip, path)
}

// GetClientIP extracts client IP from request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// Get first IP from comma-separated list
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}
	
	// Check CF-Connecting-IP (Cloudflare)
	cfip := r.Header.Get("CF-Connecting-IP")
	if cfip != "" {
		return cfip
	}
	
	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Helper function for max (Go 1.21+ has built-in max)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Predefined configurations for common use cases
var (
	// StrictConfig for high-security endpoints
	StrictConfig = Config{
		Requests:       10,
		Window:         time.Minute,
		KeyBuilder:     UserKeyBuilder,
		IncludeHeaders: true,
		Message:        "Rate limit exceeded for this endpoint",
		StatusCode:     http.StatusTooManyRequests,
	}
	
	// APIConfig for general API endpoints
	APIConfig = Config{
		Requests:       1000,
		Window:         time.Hour,
		KeyBuilder:     DefaultKeyBuilder,
		SkipPaths:      []string{"/health", "/health/*", "/metrics"},
		IncludeHeaders: true,
		Message:        "API rate limit exceeded",
		StatusCode:     http.StatusTooManyRequests,
	}
	
	// LoginConfig for authentication endpoints
	LoginConfig = Config{
		Requests:       5,
		Window:         15 * time.Minute,
		KeyBuilder:     DefaultKeyBuilder,
		IncludeHeaders: true,
		Message:        "Too many login attempts",
		StatusCode:     http.StatusTooManyRequests,
	}
)