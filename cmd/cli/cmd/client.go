package cmd

import (
	"net/url"

	"github.com/theopenlane/openlane-cloud/internal/client"
)

// SetupClient will setup the openlane cloud client
func SetupClient(host string) (client.Client, error) {
	config := client.NewDefaultConfig()

	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	return client.New(config, client.WithBaseURL(baseURL))
}
