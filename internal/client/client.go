package client

import (
	"context"

	"github.com/theopenlane/httpsling"

	"github.com/theopenlane/openlane-cloud/internal/v1/models"
)

// Client is the interface that wraps the openlane cloud API REST client methods
type Client interface {
	// OrganizationCreate creates an organizational hierarchy for a organization
	OrganizationCreate(context.Context, *models.OrganizationRequest) (*models.OrganizationReply, error)
}

// NewWithDefaults creates a new API v1 client with default configuration
func NewWithDefaults() (Client, error) {
	conf := NewDefaultConfig()

	return New(conf)
}

// New creates a new API v1 client that implements the Client interface
func New(config Config, opts ...Option) (_ Client, err error) {
	c := &APIv1{
		Config: config,
	}

	// create the HTTP sling client if it is not set
	if c.Requester == nil {
		c.Requester, err = httpsling.New()
		if err != nil {
			return nil, err
		}
	}

	// apply the options to the client
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// APIv1 implements the Client interface and provides methods to interact with the openlane cloud API
type APIv1 struct {
	// Config is the configuration for the APIv1 client
	Config Config
	// HTTPSlingClient is the HTTP client for the APIv1 client
	Requester *httpsling.Requester
}

// Ensure the APIv1 implements the Client interface
var _ Client = &APIv1{}

// OrganizationCreate creates an organizational hierarchy for a new organization based on the name, environment(s), bucket(s), and
// relationship(s) provided in the request
func (c *APIv1) OrganizationCreate(ctx context.Context, in *models.OrganizationRequest) (out *models.OrganizationReply, err error) {
	resp, err := c.Requester.ReceiveWithContext(ctx, &out,
		httpsling.Post(v1Path("organization")),
		httpsling.Body(in))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if !httpsling.IsSuccess(resp) {
		return nil, newRequestError(resp.StatusCode, out.Error)
	}

	return out, nil
}

func v1Path(path string) string {
	return "/v1/" + path
}
