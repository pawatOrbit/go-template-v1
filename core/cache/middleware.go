package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/yourorg/go-api-template/core/logger"
)

// CacheMiddlewareConfig holds cache middleware configuration
type CacheMiddlewareConfig struct {
	// DefaultTTL is the default cache expiration time
	DefaultTTL time.Duration
	// SkipCache functions that determine whether to skip caching
	SkipCache []func(*http.Request) bool
	// CacheKeyBuilder builds cache keys from requests
	CacheKeyBuilder func(*http.Request) string
	// OnlyMethods specifies which HTTP methods to cache (default: GET)
	OnlyMethods []string
	// SkipPaths are paths to skip caching
	SkipPaths []string
}

// DefaultCacheMiddlewareConfig returns a default cache middleware configuration
func DefaultCacheMiddlewareConfig() CacheMiddlewareConfig {
	return CacheMiddlewareConfig{
		DefaultTTL:      5 * time.Minute,
		OnlyMethods:     []string{"GET"},
		SkipPaths:       []string{"/health", "/health/*", "/metrics"},
		CacheKeyBuilder: DefaultCacheKeyBuilder,
		SkipCache: []func(*http.Request) bool{
			SkipAuthenticatedRequests,
		},
	}
}

// CacheMiddleware creates an HTTP middleware for caching responses
func CacheMiddleware(cacheService CacheService, config CacheMiddlewareConfig) func(http.Handler) http.Handler {
	// Set defaults if not provided
	if len(config.OnlyMethods) == 0 {
		config.OnlyMethods = []string{"GET"}
	}
	if config.DefaultTTL == 0 {
		config.DefaultTTL = 5 * time.Minute
	}
	if config.CacheKeyBuilder == nil {
		config.CacheKeyBuilder = DefaultCacheKeyBuilder
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if we should cache this request
			if !shouldCache(r, config) {
				next.ServeHTTP(w, r)
				return
			}

			// Build cache key
			cacheKey := config.CacheKeyBuilder(r)
			
			// Try to get from cache first
			ctx := r.Context()
			cachedResponse, err := getCachedResponse(ctx, cacheService, cacheKey)
			if err == nil {
				// Cache hit - serve cached response
				serveCachedResponse(w, cachedResponse)
				logger.Slog.InfoContext(ctx, "Cache hit", "key", cacheKey)
				return
			}

			// Cache miss - capture response and cache it
			responseCapture := &responseCapture{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:          make([]byte, 0),
			}

			next.ServeHTTP(responseCapture, r)

			// Cache the response if it's successful
			if responseCapture.statusCode >= 200 && responseCapture.statusCode < 300 {
				cached := &cachedResponseData{
					StatusCode: responseCapture.statusCode,
					Headers:    w.Header(),
					Body:       responseCapture.body,
					Timestamp:  time.Now(),
				}

				err := cacheResponse(ctx, cacheService, cacheKey, cached, config.DefaultTTL)
				if err != nil {
					logger.Slog.ErrorContext(ctx, "Failed to cache response", "key", cacheKey, "error", err.Error())
				} else {
					logger.Slog.InfoContext(ctx, "Response cached", "key", cacheKey, "ttl", config.DefaultTTL)
				}
			}
		})
	}
}

// shouldCache determines if a request should be cached
func shouldCache(r *http.Request, config CacheMiddlewareConfig) bool {
	// Check HTTP method
	methodAllowed := false
	for _, method := range config.OnlyMethods {
		if r.Method == method {
			methodAllowed = true
			break
		}
	}
	if !methodAllowed {
		return false
	}

	// Check skip paths
	for _, path := range config.SkipPaths {
		if path == r.URL.Path || (strings.HasSuffix(path, "/*") && strings.HasPrefix(r.URL.Path, strings.TrimSuffix(path, "/*"))) {
			return false
		}
	}

	// Check custom skip functions
	for _, skipFunc := range config.SkipCache {
		if skipFunc(r) {
			return false
		}
	}

	return true
}

// DefaultCacheKeyBuilder builds a cache key from the request
func DefaultCacheKeyBuilder(r *http.Request) string {
	// Create hash from method + path + query parameters
	h := md5.New()
	h.Write([]byte(r.Method))
	h.Write([]byte(r.URL.Path))
	h.Write([]byte(r.URL.RawQuery))
	
	hash := hex.EncodeToString(h.Sum(nil))
	return BuildCacheKey("http", r.Method, r.URL.Path, hash[:8])
}

// SkipAuthenticatedRequests skips caching for requests with Authorization header
func SkipAuthenticatedRequests(r *http.Request) bool {
	return r.Header.Get("Authorization") != ""
}

// SkipQueryParams skips caching for requests with query parameters
func SkipQueryParams(r *http.Request) bool {
	return r.URL.RawQuery != ""
}

// cachedResponseData represents a cached HTTP response
type cachedResponseData struct {
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
	Timestamp  time.Time   `json:"timestamp"`
}

// responseCapture captures HTTP response data
type responseCapture struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rc *responseCapture) WriteHeader(statusCode int) {
	rc.statusCode = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}

func (rc *responseCapture) Write(data []byte) (int, error) {
	rc.body = append(rc.body, data...)
	return rc.ResponseWriter.Write(data)
}

// getCachedResponse retrieves a cached response
func getCachedResponse(ctx context.Context, cacheService CacheService, key string) (*cachedResponseData, error) {
	var cached cachedResponseData
	err := cacheService.GetJSON(ctx, key, &cached)
	if err != nil {
		return nil, err
	}
	return &cached, nil
}

// cacheResponse stores a response in cache
func cacheResponse(ctx context.Context, cacheService CacheService, key string, response *cachedResponseData, ttl time.Duration) error {
	return cacheService.SetJSON(ctx, key, response, ttl)
}

// serveCachedResponse serves a cached response
func serveCachedResponse(w http.ResponseWriter, cached *cachedResponseData) {
	// Set headers
	for key, values := range cached.Headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	
	// Add cache headers
	w.Header().Set("X-Cache", "HIT")
	w.Header().Set("X-Cache-Time", cached.Timestamp.Format(time.RFC3339))
	
	// Write status and body
	w.WriteHeader(cached.StatusCode)
	w.Write(cached.Body)
}

// CacheInvalidateMiddleware provides cache invalidation functionality
func CacheInvalidateMiddleware(cacheService CacheService, patterns ...string) func(http.Handler) http.Handler {
	if len(patterns) == 0 {
		patterns = []string{"http:*"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only invalidate on write operations
			if r.Method != "POST" && r.Method != "PUT" && r.Method != "PATCH" && r.Method != "DELETE" {
				next.ServeHTTP(w, r)
				return
			}

			// Capture response to check if operation was successful
			responseCapture := &responseCapture{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(responseCapture, r)

			// Only invalidate on successful operations
			if responseCapture.statusCode >= 200 && responseCapture.statusCode < 300 {
				ctx := r.Context()
				for _, pattern := range patterns {
					keys, err := cacheService.Keys(ctx, pattern)
					if err != nil {
						logger.Slog.ErrorContext(ctx, "Failed to get cache keys for invalidation", "pattern", pattern, "error", err.Error())
						continue
					}

					if len(keys) > 0 {
						err = cacheService.Delete(ctx, keys...)
						if err != nil {
							logger.Slog.ErrorContext(ctx, "Failed to invalidate cache", "keys", keys, "error", err.Error())
						} else {
							logger.Slog.InfoContext(ctx, "Cache invalidated", "pattern", pattern, "keys_count", len(keys))
						}
					}
				}
			}
		})
	}
}

// Cache decorator functions for manual caching in services
func CacheGet[T any](ctx context.Context, cacheService CacheService, key string, dest *T) error {
	return cacheService.GetJSON(ctx, key, dest)
}

func CacheSet[T any](ctx context.Context, cacheService CacheService, key string, value T, ttl time.Duration) error {
	return cacheService.SetJSON(ctx, key, value, ttl)
}

func CacheGetOrSet[T any](ctx context.Context, cacheService CacheService, key string, ttl time.Duration, fetchFunc func() (T, error)) (T, error) {
	var result T
	
	// Try to get from cache first
	err := cacheService.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}
	
	// Cache miss - fetch data
	result, err = fetchFunc()
	if err != nil {
		return result, err
	}
	
	// Cache the result
	cacheErr := cacheService.SetJSON(ctx, key, result, ttl)
	if cacheErr != nil {
		logger.Slog.ErrorContext(ctx, "Failed to cache result", "key", key, "error", cacheErr.Error())
	}
	
	return result, nil
}