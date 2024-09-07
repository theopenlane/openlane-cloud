package cmd

import (
	"log"

	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const appName = "openlane-cloud"

var (
	logger *zap.SugaredLogger
	k      *koanf.Koanf
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
		log.Fatalf("error loading config: %v", err)
	}

	setupLogging()
}

func initCmdFlags(cmd *cobra.Command) error {
	return k.Load(posflag.Provider(cmd.Flags(), k.Delim(), k), nil)
}

func setupLogging() {
	cfg := zap.NewProductionConfig()
	if k.Bool("pretty") {
		cfg = zap.NewDevelopmentConfig()
	}

	if k.Bool("debug") {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = l.Sugar().With("app", appName)
	defer logger.Sync() //nolint:errcheck
}
