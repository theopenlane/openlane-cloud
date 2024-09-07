package client

import (
	"net/url"

	"github.com/theopenlane/httpsling"
)

// ClientOption allows us to configure the APIv1 client when it is created
type ClientOption func(c *APIv1) error

// WithBaseURL sets the base URL for the APIv1 client
func WithBaseURL(baseURL *url.URL) ClientOption {
	return func(c *APIv1) error {
		// Set the base URL for the APIv1 client
		c.Config.BaseURL = baseURL

		// Set the base URL for the HTTPSling client
		return c.Requester.Apply(httpsling.URL(baseURL.String()))
	}
}
