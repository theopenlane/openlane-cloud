package openlane

import (
	"github.com/theopenlane/core/pkg/openlaneclient"
)

// NewDefaultClient creates a new openlane client using the default configuration variables
func NewDefaultClient() (*openlaneclient.OpenlaneClient, error) {
	config, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	return config.createClient()
}

// NewClient creates a new openlane client using the provided configuration variables
func (c *Config) NewClient() (*openlaneclient.OpenlaneClient, error) {
	return c.createClient()
}

// CreateOpenlaneClient creates a new openlane client using the OPENLANE_TOKEN configuration variable
func (c *Config) createClient() (*openlaneclient.OpenlaneClient, error) {
	if c.Token == "" {
		return nil, ErrAPITokenMissing
	}

	config := openlaneclient.NewDefaultConfig()

	opt := openlaneclient.WithCredentials(openlaneclient.Authorization{
		BearerToken: c.Token})

	return openlaneclient.New(config, opt)
}
