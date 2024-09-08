package cmd

import (
	"context"
	"net/url"

	"github.com/theopenlane/openlane-cloud/internal/client"
)

// SetupClient will setup the openlane cloud client
func SetupClient(ctx context.Context, host string) (client.Client, error) {
	config := client.NewDefaultConfig()

	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return client.New(config, client.WithBaseURL(baseURL))
}
