package core_config

import (
	"github.com/yourorg/go-api-template/core/cache"
	"github.com/yourorg/go-api-template/core/pgdb"
)

type Config struct {
	Env        string         `mapstructure:"env"`
	RestServer RestServer     `mapstructure:"restServer"`
	CORS       CORS           `mapstructure:"cors"`
	Postgres   pgdb.Postgres  `mapstructure:"postgres"`
	LMStudio   LMStudioConfig `mapstructure:"lmStudio"`
	Auth       AuthConfig     `mapstructure:"auth"`
	Redis      cache.RedisConfig `mapstructure:"redis"`
	RateLimit  RateLimitConfig `mapstructure:"rateLimit"`
}

type CORS struct {
	AllowedMethods []string `mapstructure:"allowedMethods"`
	AllowedHeaders []string `mapstructure:"allowedHeaders"`
	AllowedOrigins []string `mapstructure:"allowedOrigins"` // Default: ["*"]
	ExposedHeaders []string `mapstructure:"exposedHeaders"`
	MaxAge         int      `mapstructure:"maxAge"` // Default: 7200 (seconds)
}

type RestServer struct {
	Port string `mapstructure:"port"`
}

type LMStudioConfig struct {
	Protocol    string  `mapstructure:"protocol"`
	BaseUrl     string  `mapstructure:"baseUrl"`
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"maxTokens"`
	EnableMock  bool    `mapstructure:"enableMock"`
}

type AuthConfig struct {
	JWTSecretKey   string   `mapstructure:"jwtSecretKey"`
	SkipAuthPaths  []string `mapstructure:"skipAuthPaths"`
	TokenDuration  string   `mapstructure:"tokenDuration"`  // e.g., "24h"
	RefreshDuration string  `mapstructure:"refreshDuration"` // e.g., "168h" (7 days)
}

type RateLimitConfig struct {
	Enabled       bool     `mapstructure:"enabled"`
	Requests      int      `mapstructure:"requests"`
	Window        string   `mapstructure:"window"`        // e.g., "1h", "15m"
	SkipPaths     []string `mapstructure:"skipPaths"`
	IncludeHeaders bool    `mapstructure:"includeHeaders"`
	Message       string   `mapstructure:"message"`
	StatusCode    int      `mapstructure:"statusCode"`
}
