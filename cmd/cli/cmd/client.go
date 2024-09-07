package cmd

import (
	"context"
	"net/url"

	"github.com/theopenlane/openlane-cloud/internal/client"
)

// SetupClient will setup the openlane cloud client
func SetupClient(ctx context.Context, host string) (client.Client, error) {
	config := client.NewDefaultConfig()

	opt, err := configureClientEndpoints(host)
	if err != nil {
		return nil, err
	}

	return client.New(config, opt)
}

// configureClientEndpoints will setup the base URL for the openlane client
func configureClientEndpoints(host string) (client.ClientOption, error) {
	if host == "" {
		return nil, nil
	}

	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return client.WithBaseURL(baseURL), nil
}
