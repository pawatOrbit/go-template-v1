/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/pawatOrbit/ai-mock-data-service/go/config"
	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/logger"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/pgdb"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/build"
	"github.com/pawatOrbit/ai-mock-data-service/go/utils/runtime"
	"github.com/spf13/cobra"
)

var RootCmdName = "main"

var rootCmd = &cobra.Command{
	Use:   "help",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().String("profile", "", "Profile for the service to run")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// To get config at runtime
func getConfigFunc() core_config.Config {
	return config.GetConfig().Config
}

func setUpLogger(validateProfile runtime.Environment) {
	logger.InitLogger(validateProfile)
}

func setUpConfig(profile runtime.Environment) {
	// Read config from config file

	runtimeCfg := runtime.RuntimeCfg{
		Microservice: build.ServiceName,
		Env:          profile,
	}

	ctx := context.Background()

	globalConfigPath, err := core_config.GetGlobalConfigFilePath(runtimeCfg)
	if err != nil {
		fmt.Println("Error getting global config file path", err.Error())
		return
	}

	slog.InfoContext(ctx, "Getting config from file", "file", globalConfigPath)
	err = config.ResolveConfigFromFile(ctx, globalConfigPath)
	if err != nil {
		fmt.Println("Error reading global config file", err.Error())
	}
}

func setUpPostgres() {
	postgresConfig := config.GetConfig().Postgres

	ctx := context.Background()
	// Check if either Read or Write host is provided
	if postgresConfig.Read.Host != "" || postgresConfig.Write.Host != "" {
		// Initialize database schema if it doesn't exist
		ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		// Ensure both Read and Write hosts are provided
		if postgresConfig.Read.Schema == "" || postgresConfig.Write.Schema == "" {
			slog.ErrorContext(ctx, "Both Read and Write schema must be set for Postgres")
		}

		// Initialize the connection pool
		slog.InfoContext(ctx, "Initializing pgxPool")
		err := pgdb.InitPgConnectionPool(ctx, postgresConfig)
		if err != nil {
			slog.ErrorContext(ctx, "Failed in pgx.InitPgConnectionPool()", "error", err)
		}
		slog.InfoContext(ctx, "pgxPool initialized")
	}
}
