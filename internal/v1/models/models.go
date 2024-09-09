package models

import (
	"github.com/mcuadros/go-defaults"
	"github.com/theopenlane/utils/rout"
)

// =========
// ORGANIZATION
// =========

// OrganizationRequest is the request object for creating a organization
type OrganizationRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Domains       []string `json:"domains,omitempty"`
	Environments  []string `json:"environments,omitempty" default:"[production,testing]"`
	Buckets       []string `json:"buckets,omitempty" default:"[assets,customers,orders,relationships,sales]"`
	Relationships []string `json:"relationships,omitempty" default:"[internal_users,marketing_subscribers,marketplaces,partners,vendors]"`
}

// OrganizationReply is the response object for creating a organization
type OrganizationReply struct {
	rout.Reply
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	Domains      []string      `json:"domains,omitempty"`
	Environments []Environment `json:"environments,omitempty"`
}

type Environment struct {
	OrgDetails
	Buckets []Bucket `json:"buckets,omitempty"`
}

type Bucket struct {
	OrgDetails
	Relations []Relationship `json:"relations,omitempty"`
}

type Relationship struct {
	OrgDetails
}

type OrgDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Validate ensures the required fields are set on the OrganizationRequest request
func (r *OrganizationRequest) Validate() error {
	// Required for all requests
	if r.Name == "" {
		return rout.MissingField("name")
	}

	// Set default values if not provided in the request
	defaultRequest := &OrganizationRequest{}
	defaults.SetDefaults(defaultRequest)

	if r.Environments == nil {
		r.Environments = defaultRequest.Environments
	}

	if r.Buckets == nil {
		r.Buckets = defaultRequest.Buckets
	}

	if r.Relationships == nil {
		r.Relationships = defaultRequest.Relationships
	}

	return nil
}

// ExampleOrganizationSuccessRequest is an example of a successful organization request for OpenAPI documentation
var ExampleOrganizationSuccessRequest = OrganizationRequest{
	Name: "MITB Inc.",
}

// ExampleOrganizationSuccessResponse is an example of a successful organization response for OpenAPI documentation
var ExampleOrganizationSuccessResponse = OrganizationReply{
	Reply: rout.Reply{Success: true},
	ID:    "1234",
	Name:  "MITB Inc.",
}
