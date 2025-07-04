package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	core_config "github.com/pawatOrbit/ai-mock-data-service/go/core/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/server"
	"github.com/pawatOrbit/ai-mock-data-service/go/utils/runtime"
	"github.com/spf13/cobra"
)

func init() {
	InitServeCommandGroup(rootCmd)

	servePreRunFunc := func(cmd *cobra.Command, _ []string) {
		profile, err := cmd.Flags().GetString("profile")
		if err != nil {
			slog.Error("Error getting profile flag", slog.Any("Error", err))
		}
		validatedProfile := runtime.ValidateProfile(profile)
		setUpLogger(validatedProfile)
		setUpConfig(validatedProfile)
		setUpPostgres()
	}

	NewServe(
		rootCmd,
		servePreRunFunc,
		getConfigFunc,
		WithHTTPServer(server.NewHttpServer),
	)
}

type ServeOpts struct {
	initHTTPServer func() (*http.Server, error)
}

func WithHTTPServer(fn func() (*http.Server, error)) ServeOptsFunc {
	return func(o *ServeOpts) {
		o.initHTTPServer = fn
	}
}

func defaultServeOpts() ServeOpts {
	return ServeOpts{}
}

type ServeOptsFunc func(*ServeOpts)

func NewServe(rootCmd *cobra.Command, preRunFunc func(cmd *cobra.Command, _ []string), getConfig func() core_config.Config, serveOpts ...ServeOptsFunc) *cobra.Command {
	command := cobra.Command{
		Use:     "serve:all-api",
		Short:   "Start REST API server",
		GroupID: "serve",
		PreRun:  preRunFunc,
		RunE: func(cmd *cobra.Command, args []string) error {
			o := defaultServeOpts()
			for _, f := range serveOpts {
				f(&o)
			}
			ctx := cmd.Context()
			cfg := getConfig()
			restPort := cfg.RestServer.Port
			localIP, _ := getLocalIP()

			if o.initHTTPServer != nil {
				restServer, err := o.initHTTPServer()
				if err != nil {
					return fmt.Errorf("failed to create REST server: %w", err)
				}
				go func() {
					slog.InfoContext(ctx, fmt.Sprintf("[REST] Starting server on port %s", restPort))
					slog.InfoContext(ctx, fmt.Sprintf("[REST] Local: http://localhost:%s", restPort))
					slog.InfoContext(ctx, fmt.Sprintf("[REST] Network: http://%s:%s", localIP, restPort))
					slog.InfoContext(ctx, "[REST] waiting for requests...")
					if err := restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						slog.ErrorContext(ctx, fmt.Sprintf("[REST] failed to serve: %s\n", err))
					}

				}()

				go func() {
					<-ctx.Done()
					gracefulShutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
					defer cancel()
					restServer.Shutdown(gracefulShutdownCtx)
				}()

			}

			<-ctx.Done()
			return nil
		},
	}

	rootCmd.AddCommand(&command)
	return &command
}
