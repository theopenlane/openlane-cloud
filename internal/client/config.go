package client

import (
	"net/url"
)

// Config is the configuration for the openlane cloud API client
type Config struct {
	// BaseURL is the base URL for the openlane API
	BaseURL *url.URL `json:"baseUrl" yaml:"base_url" default:"http://localhost:17610"`
}

// NewDefaultConfig returns a new default configuration for the openlane cloud API client
func NewDefaultConfig() Config {
	return defaultClientConfig
}

var defaultClientConfig = Config{
	BaseURL: &url.URL{
		Scheme: "http",
		Host:   "localhost:17610",
	},
}
