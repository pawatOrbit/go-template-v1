package core_config

import "github.com/pawatOrbit/ai-mock-data-service/go/core/pgdb"

type Config struct {
	Env        string         `mapstructure:"env"`
	RestServer RestServer     `mapstructure:"restServer"`
	CORS       CORS           `mapstructure:"cors"`
	Postgres   pgdb.Postgres  `mapstructure:"postgres"`
	LMStudio   LMStudioConfig `mapstructure:"lmStudio"`
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
