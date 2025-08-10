package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourorg/go-api-template/core/logger"
)

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	Database     int           `mapstructure:"database"`
	MaxRetries   int           `mapstructure:"maxRetries"`
	PoolSize     int           `mapstructure:"poolSize"`
	MinIdleConns int           `mapstructure:"minIdleConns"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

// CacheService provides caching functionality
type CacheService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
	FlushDB(ctx context.Context) error
	Close() error
	Ping(ctx context.Context) error
	GetClient() *redis.Client
}

type redisService struct {
	client *redis.Client
	config RedisConfig
}

var (
	redisInstance CacheService
	redisOnce     sync.Once
)

// NewRedisService creates a new Redis service
func NewRedisService(config RedisConfig) CacheService {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.Database,
		MaxRetries:   config.MaxRetries,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	return &redisService{
		client: client,
		config: config,
	}
}

// GetRedisService returns singleton Redis service instance
func GetRedisService() CacheService {
	return redisInstance
}

// InitRedisService initializes the global Redis service
func InitRedisService(config RedisConfig) error {
	var err error
	redisOnce.Do(func() {
		redisInstance = NewRedisService(config)
		
		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		err = redisInstance.Ping(ctx)
		if err != nil {
			logger.Slog.Error("Failed to connect to Redis", "error", err.Error())
		} else {
			logger.Slog.Info("Redis connection established successfully")
		}
	})
	
	return err
}

// Get retrieves a string value from cache
func (r *redisService) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheKeyNotFound
	}
	if err != nil {
		logger.Slog.Error("Redis GET error", "key", key, "error", err.Error())
		return "", fmt.Errorf("redis get error: %w", err)
	}
	return result, nil
}

// Set stores a value in cache with expiration
func (r *redisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		logger.Slog.Error("Redis SET error", "key", key, "error", err.Error())
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

// GetJSON retrieves and unmarshals JSON from cache
func (r *redisService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	result, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		logger.Slog.Error("JSON unmarshal error", "key", key, "error", err.Error())
		return fmt.Errorf("json unmarshal error: %w", err)
	}
	
	return nil
}

// SetJSON marshals and stores JSON in cache
func (r *redisService) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		logger.Slog.Error("JSON marshal error", "key", key, "error", err.Error())
		return fmt.Errorf("json marshal error: %w", err)
	}
	
	return r.Set(ctx, key, jsonData, expiration)
}

// Delete removes one or more keys from cache
func (r *redisService) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		logger.Slog.Error("Redis DELETE error", "keys", keys, "error", err.Error())
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}

// Exists checks if a key exists in cache
func (r *redisService) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		logger.Slog.Error("Redis EXISTS error", "key", key, "error", err.Error())
		return false, fmt.Errorf("redis exists error: %w", err)
	}
	return result > 0, nil
}

// Expire sets expiration time for a key
func (r *redisService) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		logger.Slog.Error("Redis EXPIRE error", "key", key, "error", err.Error())
		return fmt.Errorf("redis expire error: %w", err)
	}
	return nil
}

// TTL returns the time to live for a key
func (r *redisService) TTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		logger.Slog.Error("Redis TTL error", "key", key, "error", err.Error())
		return 0, fmt.Errorf("redis ttl error: %w", err)
	}
	return result, nil
}

// Keys returns all keys matching pattern
func (r *redisService) Keys(ctx context.Context, pattern string) ([]string, error) {
	result, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		logger.Slog.Error("Redis KEYS error", "pattern", pattern, "error", err.Error())
		return nil, fmt.Errorf("redis keys error: %w", err)
	}
	return result, nil
}

// FlushDB removes all keys from current database
func (r *redisService) FlushDB(ctx context.Context) error {
	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		logger.Slog.Error("Redis FLUSHDB error", "error", err.Error())
		return fmt.Errorf("redis flushdb error: %w", err)
	}
	return nil
}

// Ping tests connectivity to Redis
func (r *redisService) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping error: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (r *redisService) Close() error {
	return r.client.Close()
}

// GetClient returns the underlying Redis client
func (r *redisService) GetClient() *redis.Client {
	return r.client
}

// Cache key helpers
func BuildCacheKey(parts ...string) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}

// Common expiration times
const (
	ExpireNever     = time.Duration(-1)
	Expire1Minute   = time.Minute
	Expire5Minutes  = 5 * time.Minute
	Expire15Minutes = 15 * time.Minute
	Expire30Minutes = 30 * time.Minute
	Expire1Hour     = time.Hour
	Expire6Hours    = 6 * time.Hour
	Expire12Hours   = 12 * time.Hour
	Expire24Hours   = 24 * time.Hour
	Expire7Days     = 7 * 24 * time.Hour
)

// Custom errors
var (
	ErrCacheKeyNotFound = fmt.Errorf("cache key not found")
)