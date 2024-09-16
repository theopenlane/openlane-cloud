// Package cmd is our cobra cli implementation
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	appName = "openlane-cloud"
)

var (
	cfgFile      string
	OutputFormat string
	Config       *koanf.Koanf
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   appName,
	Short: "the openlane-cloud cli",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initConfiguration(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	Config = koanf.New(".") // Create a new koanf instance.

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+appName+".yaml)")
	RootCmd.PersistentFlags().String("host", "", "openlane cloud api url, if not set the config default will be used")
	RootCmd.PersistentFlags().String("openlanehost", "", "openlane api url, if not set the config default will be used")

	// Auth flags
	RootCmd.Flags().String("token", "", "openlane api token")

	// Logging flags
	RootCmd.PersistentFlags().Bool("debug", false, "enable debug logging")
	RootCmd.PersistentFlags().Bool("pretty", false, "enable pretty (human readable) logging output")

	// Output flags
	RootCmd.PersistentFlags().StringVarP(&OutputFormat, "format", "z", "table", "output format (json, table)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// load the flags to ensure we know the correct config file path
	initConfiguration(RootCmd)

	// load the config file and env vars
	loadConfigFile()

	// reload because flags and env vars take precedence over file
	initConfiguration(RootCmd)

	// setup logging configuration
	setupLogging()
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
	if Config.Bool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)

		buildInfo, _ := debug.ReadBuildInfo()

		log.Logger = log.Logger.With().
			Caller().
			Int("pid", os.Getpid()).
			Str("go_version", buildInfo.GoVersion).Logger()
	}

	// pretty logging for development
	if Config.Bool("pretty") {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			FormatCaller: func(i interface{}) string {
				return filepath.Base(fmt.Sprintf("%s", i))
			},
		})
	}
}

// initConfiguration loads the configuration from the command flags of the given cobra command
func initConfiguration(cmd *cobra.Command) {
	loadEnvVars()

	loadFlags(cmd)
}

func loadConfigFile() {
	if cfgFile == "" {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		cfgFile = filepath.Join(home, "."+appName+".yaml")
	}

	// If the config file does not exist, do nothing
	if _, err := os.Stat(cfgFile); errors.Is(err, os.ErrNotExist) {
		return
	}

	err := Config.Load(file.Provider(cfgFile), yaml.Parser())

	cobra.CheckErr(err)
}

func loadEnvVars() {
	err := Config.Load(env.ProviderWithValue("OPENLANECLOUD_", ".", func(s string, v string) (string, interface{}) {
		key := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "OPENLANECLOUD_")), "_", ".")

		if strings.Contains(v, ",") {
			return key, strings.Split(v, ",")
		}

		return key, v
	}), nil)

	cobra.CheckErr(err)
}

func loadFlags(cmd *cobra.Command) {
	err := Config.Load(posflag.Provider(cmd.Flags(), Config.Delim(), Config), nil)

	cobra.CheckErr(err)
}
