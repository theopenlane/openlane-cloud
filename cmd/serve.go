package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/theopenlane/beacon/otelx"
	"go.uber.org/zap"

	"github.com/theopenlane/openlane-cloud/internal/httpserve/config"
	"github.com/theopenlane/openlane-cloud/internal/httpserve/server"
	"github.com/theopenlane/openlane-cloud/internal/httpserve/serveropts"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the openlane cloud API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().String("config", "./config/.config.yaml", "config file location")
}

func serve(ctx context.Context) error {
	serverOpts := []serveropts.ServerOption{}
	serverOpts = append(serverOpts,
		serveropts.WithConfigProvider(&config.ConfigProviderWithRefresh{}),
		serveropts.WithLogger(logger),
		serveropts.WithOpenlaneClient(),
		serveropts.WithHTTPS(),
		serveropts.WithMiddleware(),
		serveropts.WithRateLimiter(),
	)

	so := serveropts.NewServerOptions(serverOpts, k.String("config"))

	if err := otelx.NewTracer(so.Config.Settings.Tracer, appName); err != nil {
		logger.Fatalw("failed to initialize tracer", "error", err)
	}

	srv := server.NewServer(so.Config, so.Config.Logger)

	if err := srv.StartEchoServer(ctx); err != nil {
		logger.Error("failed to run server", zap.Error(err))
	}

	return nil
}
