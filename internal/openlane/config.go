package openlane

import "github.com/mcuadros/go-defaults"

// Config is the configuration for the openlane api
type Config struct {
	// Token is the token to use for the openlane client
	Token string `json:"token" koanf:"token" default:""`
}

// NewDefaultConfig returns a new Config with default values
func NewDefaultConfig() (*Config, error) {
	// Set default values
	conf := &Config{}
	defaults.SetDefaults(conf)

	return conf, nil
}
