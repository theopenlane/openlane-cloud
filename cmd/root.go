package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const appName = "openlane-cloud"

var (
	k *koanf.Koanf
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "cli for interacting with the openlane cloud server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initCmdFlags(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	k = koanf.New(".") // Create a new koanf instance.

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().Bool("pretty", false, "enable pretty (human readable) logging output")
	rootCmd.PersistentFlags().Bool("debug", false, "debug logging output")
}

// initConfig reads in flags set for server startup
// all other configuration is done by the server with koanf
// refer to the README.md for more information
func initConfig() {
	// Load config from flags, including defaults
	if err := initCmdFlags(rootCmd); err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	setupLogging()
}

func initCmdFlags(cmd *cobra.Command) error {
	return k.Load(posflag.Provider(cmd.Flags(), k.Delim(), k), nil)
}

// setupLogging sets up the logging defaults for the application
func setupLogging() {
	// setup logging with time and app name
	log.Logger = zerolog.New(os.Stderr).
		With().Timestamp().
		Logger().
		With().Str("app", appName).
		Logger()

	// set the log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// set the log level to debug if the debug flag is set and add additional information
	if k.Bool("pretty") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		buildInfo, _ := debug.ReadBuildInfo()

		log.Logger = log.Logger.With().
			Caller().
			Int("pid", os.Getpid()).
			Str("go_version", buildInfo.GoVersion).Logger()
	}

	// pretty logging for development
	if k.Bool("pretty") {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatCaller: func(i interface{}) string {
				return filepath.Base(fmt.Sprintf("%s", i))
			},
		})
	}
}
