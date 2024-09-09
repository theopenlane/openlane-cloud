package seed

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/url"
	"time"

	"github.com/theopenlane/core/pkg/enums"
	"github.com/theopenlane/core/pkg/models"
	"github.com/theopenlane/core/pkg/openlaneclient"
)

// Config represents provides the openlane client and configuration for the seed client
type Client struct {
	*openlaneclient.OpenlaneClient
	config *Config
}

// NewDefaultClient creates a new openlane client using the default configuration variables
func NewDefaultClient() (*Client, error) {
	config, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	client, err := config.newOpenlaneClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		OpenlaneClient: client,
		config:         config,
	}, nil
}

// NewClient creates a new openlane client using the provided configuration variables
func (c *Config) NewClient() (*Client, error) {
	client, err := c.newOpenlaneClient()
	if err != nil {
		return nil, err
	}

	return &Client{
		OpenlaneClient: client,
		config:         c,
	}, nil
}

func (c *Config) newOpenlaneClient() (*openlaneclient.OpenlaneClient, error) {
	config := openlaneclient.NewDefaultConfig()

	var err error

	// if the openlane host is set, use it, otherwise use the default from the config
	if c.OpenlaneHost != "" {
		config.BaseURL, err = url.Parse(c.OpenlaneHost)
		if err != nil {
			return nil, err
		}
	}

	opts := []openlaneclient.ClientOption{openlaneclient.WithCredentials(openlaneclient.Authorization{
		BearerToken: c.Token}),
		openlaneclient.WithBaseURL(config.BaseURL),
	}

	return openlaneclient.New(config, opts...)
}

// AuthorizeOrganizationOnPAT authorizes the organization id on the personal access token id
func (c *Client) AuthorizeOrganizationOnPAT(ctx context.Context, orgID, patID string) error {
	input := openlaneclient.UpdatePersonalAccessTokenInput{
		AddOrganizationIDs: []string{orgID},
	}

	if _, err := c.UpdatePersonalAccessToken(ctx, patID, input); err != nil {
		return err
	}

	return nil
}

// GenerateSeedAPIToken generates an API token for the organization id to use for seeding
// and authenticates the client with the token for future requests
func (c *Client) GenerateSeedAPIToken(ctx context.Context, orgID string) error {
	expiresAt := time.Now().Add(time.Hour)

	input := openlaneclient.CreateAPITokenInput{
		Name:      fmt.Sprintf("seed token %s", orgID),
		OwnerID:   &orgID,
		ExpiresAt: &expiresAt, // expires in 1 hour
		Scopes:    []string{"read", "write"},
	}

	token, err := c.CreateAPIToken(ctx, input)
	if err != nil {
		return err
	}

	// Use the token to authenticate
	c.config.Token = token.CreateAPIToken.APIToken.Token

	// create a new client with the new token
	c.OpenlaneClient, err = c.config.newOpenlaneClient()
	if err != nil {
		return err
	}

	return nil
}

// LoadOrganizations loads the organizations from the organizations.csv file
func (c *Client) LoadOrganizations(ctx context.Context) (string, error) {
	file := c.config.getOrgFilePath()

	upload, err := loadCSVFile(file)
	if err != nil {
		return "", err
	}

	org, err := c.CreateBulkCSVOrganization(ctx, upload)
	if err != nil {
		return "", err
	}

	// get the first org, this is the root org
	if len(org.CreateBulkCSVOrganization.Organizations) > 0 {
		return org.CreateBulkCSVOrganization.Organizations[0].ID, nil
	}

	return "", nil
}

// LoadGroups loads the groups from the groups.csv file
func (c *Client) LoadGroups(ctx context.Context) error {
	file := c.config.getGroupFilePath()

	upload, err := loadCSVFile(file)
	if err != nil {
		return err
	}

	if _, err := c.CreateBulkCSVGroup(ctx, upload); err != nil {
		return err
	}

	return nil
}

// LoadInvites loads the invites from the invites.csv file
func (c *Client) LoadInvites(ctx context.Context) error {
	file := c.config.getInviteFilePath()

	upload, err := loadCSVFile(file)
	if err != nil {
		return err
	}

	if _, err := c.CreateBulkCSVInvite(ctx, upload); err != nil {
		return err
	}

	return nil
}

// LoadOrgMembers loads orgs members from the user ids provided
func (c *Client) LoadOrgMembers(ctx context.Context, userIDs []string) error {
	for _, userID := range userIDs {
		_, err := c.AddUserToOrgWithRole(ctx, openlaneclient.CreateOrgMembershipInput{
			UserID: userID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadSubscribers loads the subscribers from the subscribers.csv file
func (c *Client) LoadSubscribers(ctx context.Context) error {
	file := c.config.getSubscriberFilePath()

	upload, err := loadCSVFile(file)
	if err != nil {
		return err
	}

	if _, err := c.CreateBulkCSVSubscriber(ctx, upload); err != nil {
		return err
	}

	return nil
}

// RegisterUsers registers the users from the users.csv file
func (c *Client) RegisterUsers(ctx context.Context) ([]string, error) {
	userIDs := []string{}

	file := c.config.getUserFilePath()

	upload, err := loadCSVFile(file)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(upload.File)
	records, _ := reader.ReadAll()

	for i, record := range records {
		// skip header row
		if i == 0 {
			continue
		}

		req := models.RegisterRequest{
			Email:     record[2],
			Password:  record[3],
			FirstName: record[0],
			LastName:  record[1],
		}

		reply, err := c.Register(ctx, &req)
		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, reply.ID)

		if record[6] == "true" {
			// sleep for a 100ms to avoid rate limiting
			time.Sleep(100 * time.Millisecond) // nolint:mnd

			// verify the user - this will probably break in the future when we stop
			// returning the token in the register response
			if _, err := c.VerifyEmail(ctx, &models.VerifyRequest{
				Token: reply.Token,
			}); err != nil {
				return nil, err
			}
		}
	}

	return userIDs, nil
}

// LoadTemplates loads the templates from the jsonschema/templates directory
func (c *Client) LoadTemplates(ctx context.Context) error {
	if !c.config.GenerateTemplates {
		return nil
	}

	tmpls, err := getTemplates(templateDirectory)
	if err != nil {
		return err
	}

	input := []*openlaneclient.CreateTemplateInput{}

	for _, t := range tmpls {
		input = append(input, &openlaneclient.CreateTemplateInput{
			Name:         t.Name,
			Jsonconfig:   t.JSONConfig,
			TemplateType: &enums.RootTemplate,
		})
	}

	if _, err := c.CreateBulkTemplate(ctx, input); err != nil {
		return err
	}

	return nil
}
