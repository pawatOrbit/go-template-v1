package config

import (
	"context"
	"log/slog"
	"sync"

	"dario.cat/mergo"
	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
	"github.com/spf13/viper"
)

var finalConfig *Config
var cfgFromFile *Config

var m sync.Mutex

type Config struct {
	core_config.Config `mapstructure:",squash"`
}

func NewConfig(cfg Config) {
	finalConfig = &cfg
}

func ResolveConfigFromFile(ctx context.Context, configPath string) error {
	m.Lock()
	defer m.Unlock()

	if finalConfig == nil {
		finalConfig = &Config{}
	}

	viper.SetConfigName(configPath)
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting config file", "error", err)
		return err
	}

	err = viper.Unmarshal(&cfgFromFile)
	if err != nil {
		slog.ErrorContext(ctx, "Error unmarshalling config file", "error", err)
		return err
	}

	err = mergo.Merge(finalConfig, cfgFromFile, mergo.WithOverride)
	if err != nil {
		slog.ErrorContext(ctx, "Error merging config from secrets", "error", err)
		return err
	}

	// ...
	return nil
}

func GetConfig() *Config {
	return finalConfig
}
