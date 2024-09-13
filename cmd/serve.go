package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/theopenlane/beacon/otelx"

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
		serveropts.WithOpenlaneClient(),
		serveropts.WithHTTPS(),
		serveropts.WithMiddleware(),
		serveropts.WithRateLimiter(),
	)

	so := serveropts.NewServerOptions(serverOpts, k.String("config"))

	if err := otelx.NewTracer(so.Config.Settings.Tracer, appName); err != nil {
		log.Fatal().Err(err).Msg("failed to initialize tracer")
	}

	srv := server.NewServer(so.Config)

	if err := srv.StartEchoServer(ctx); err != nil {
		log.Error().Err(err).Msg("failed to run server")
	}

	return nil
}
