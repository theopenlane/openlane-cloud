package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rs/zerolog/log"
	"github.com/theopenlane/core/pkg/openlaneclient"
	echo "github.com/theopenlane/echox"
	"github.com/theopenlane/utils/rout"

	"github.com/theopenlane/openlane-cloud/internal/v1/models"
)

const (
	relationBucketName = "relationships"
)

// OrganizationHandler is the handler for the organization endpoint
func (h *Handler) OrganizationHandler(ctx echo.Context) error {
	var in models.OrganizationRequest
	if err := ctx.Bind(&in); err != nil {
		return h.InvalidInput(ctx, err)
	}

	if err := in.Validate(); err != nil {
		return h.InvalidInput(ctx, err)
	}

	log.Debug().Str("name", in.Name).
		Strs("environments", in.Environments).
		Strs("buckets", in.Buckets).
		Strs("relationships", in.Relationships).
		Msg("creating organization")

	// create root organization
	rootOrgName := in.Name
	input := openlaneclient.CreateOrganizationInput{
		Name:        strings.ToLower(rootOrgName),
		DisplayName: &rootOrgName,
	}

	if in.Description != "" {
		input.Description = &in.Description
	}

	if len(in.Domains) > 0 {
		input.CreateOrgSettings = &openlaneclient.CreateOrganizationSettingInput{
			Domains: in.Domains,
		}
	}

	ws, err := h.OpenlaneClient.CreateOrganization(ctx.Request().Context(), input, nil)
	if err != nil {
		return h.BadRequest(ctx, err)
	}

	organization := ws.CreateOrganization.Organization

	out := models.OrganizationReply{
		Reply:       rout.Reply{Success: true},
		ID:          organization.ID,
		Name:        organization.DisplayName,
		Description: *organization.Description,
		Domains:     organization.Setting.Domains,
	}

	// create environments
	envOrgs, err := h.createEnvironments(ctx.Request().Context(), organization.ID, in.Environments, input)
	if err != nil {
		return h.BadRequest(ctx, err)
	}

	// for each environment, create buckets
	for i, envOrg := range envOrgs {
		// add environments to the response
		out.Environments = append(out.Environments, models.Environment{
			OrgDetails: models.OrgDetails{
				ID:   envOrg.ID,
				Name: envOrg.DisplayName,
			},
			Buckets: []models.Bucket{},
		})

		// create buckets
		bucketOrgs, err := h.createBuckets(ctx.Request().Context(), envOrg.ID, envOrg.DisplayName, in.Buckets, input)
		if err != nil {
			return h.BadRequest(ctx, err)
		}

		for j, bucketOrg := range bucketOrgs {
			// add buckets to the response
			out.Environments[i].Buckets = append(out.Environments[i].Buckets, models.Bucket{
				OrgDetails: models.OrgDetails{
					ID:   bucketOrg.ID,
					Name: bucketOrg.DisplayName,
				},
			})

			// create relationships under the relationships bucket
			if bucketOrg.DisplayName == relationBucketName {
				relationshipOrgs, err := h.createRelationships(ctx.Request().Context(), bucketOrg.ID, envOrg.DisplayName, in.Relationships, input)
				if err != nil {
					return h.BadRequest(ctx, err)
				}

				out.Environments[i].Buckets[j].Relations = []models.Relationship{}

				for _, relationshipOrg := range relationshipOrgs {
					// add relationships to the response
					out.Environments[i].Buckets[j].Relations = append(out.Environments[i].Buckets[j].Relations, models.Relationship{
						OrgDetails: models.OrgDetails{
							ID:   relationshipOrg.ID,
							Name: relationshipOrg.DisplayName,
						},
					})
				}
			}
		}
	}

	return h.Success(ctx, out)
}

// BindOrganizationHandler is used to bind the organization endpoint to the OpenAPI schema
func (h *Handler) BindOrganizationHandler() *openapi3.Operation {
	register := openapi3.NewOperation()
	register.Description = "Organization creates an opinionated organization hierarchy for the new organization"
	register.OperationID = "OrganizationHandler"
	register.Security = &openapi3.SecurityRequirements{}

	h.AddRequestBody("OrganizationRequest", models.ExampleOrganizationSuccessRequest, register)
	h.AddResponse("OrganizationReply", "success", models.ExampleOrganizationSuccessResponse, register, http.StatusOK)
	register.AddResponse(http.StatusInternalServerError, internalServerError())
	register.AddResponse(http.StatusBadRequest, badRequest())

	return register
}

// createChildOrganizations creates the child organizations for the organization
func (h *Handler) createChildOrganizations(ctx context.Context, namePrefix, parentOrgID string, childNames, additionalTags []string) ([]openlaneclient.CreateOrganization_CreateOrganization_Organization, error) {
	var orgs []openlaneclient.CreateOrganization_CreateOrganization_Organization

	for _, childName := range childNames {
		log.Debug().Str("childName", childName).Msg("creating child organization")

		orgName := childName

		input := openlaneclient.CreateOrganizationInput{
			Name:        strings.ToLower(fmt.Sprintf("%s.%s", namePrefix, orgName)),
			DisplayName: &orgName,
			ParentID:    &parentOrgID,
			Tags:        append(additionalTags, orgName),
		}

		// create child organization
		o, err := h.OpenlaneClient.CreateOrganization(ctx, input, nil)
		if err != nil {
			return nil, err
		}

		orgs = append(orgs, o.CreateOrganization.Organization)
	}

	return orgs, nil
}

// createEnvironments creates the environments for the organization
func (h *Handler) createEnvironments(ctx context.Context, rootOrgID string, environments []string, input openlaneclient.CreateOrganizationInput) ([]openlaneclient.CreateOrganization_CreateOrganization_Organization, error) {
	return h.createChildOrganizations(ctx, input.Name, rootOrgID, environments, []string{})
}

// createBuckets creates the buckets for the organization for each environment
func (h *Handler) createBuckets(ctx context.Context, envOrgID, environment string, buckets []string, input openlaneclient.CreateOrganizationInput) ([]openlaneclient.CreateOrganization_CreateOrganization_Organization, error) {
	return h.createChildOrganizations(ctx, fmt.Sprintf("%s.%s", input.Name, environment), envOrgID, buckets, []string{environment})
}

// createRelationships creates the relationships for the organization for each environment
func (h *Handler) createRelationships(ctx context.Context, relationshipOrgID, environment string, relationships []string, input openlaneclient.CreateOrganizationInput) ([]openlaneclient.CreateOrganization_CreateOrganization_Organization, error) {
	return h.createChildOrganizations(ctx, fmt.Sprintf("%s.%s.%s", input.Name, environment, "relationships"), relationshipOrgID, relationships, []string{environment, "relationships"})
}
