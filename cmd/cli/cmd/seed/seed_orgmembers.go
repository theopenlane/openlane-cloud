package seed

import (
	"context"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/theopenlane/core/pkg/openlaneclient"

	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"
	"github.com/theopenlane/openlane-cloud/internal/seed"
)

var seedOrgMembersCmd = &cobra.Command{
	Use:   "org-members",
	Short: "add users to an existing seeded organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return initOrgMemberData(cmd.Context())
	},
}

func init() {
	seedCmd.AddCommand(seedOrgMembersCmd)

	seedOrgMembersCmd.Flags().StringP("organization-id", "o", "", "organization ID to add users to")
	seedOrgMembersCmd.Flags().Int("users", defaultObjectCount, "number of users to generate")
	seedOrgMembersCmd.Flags().StringP("patid", "t", "", "personal access token ID to authorize the organization")
}

func initOrgMemberData(ctx context.Context) error {
	orgID := cmd.Config.String("organization-id")
	if orgID == "" {
		cobra.CheckErr("Organization ID not provided")
	}

	c, err := newSeedClient()
	cobra.CheckErr(err)

	// generate users in csv
	config, err := seed.NewDefaultConfig()
	cobra.CheckErr(err)

	config.NumUsers = cmd.Config.Int("users")

	err = config.GenerateUserData()
	cobra.CheckErr(err)

	bar := progressbar.NewOptions(100, //nolint:mnd
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15), //nolint:mnd
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetDescription("[light_green]>[reset] creating seeded org members..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[light_green]=[reset]",
			SaucerHead:    "[light_green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	defer bar.Exit() //nolint:errcheck

	cmd.BarAdd(bar, 10) //nolint:mnd

	bar.Describe("[light_green]>[reset] registering users...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	userIDs, err := c.RegisterUsers(ctx)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] creating org members...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	err = c.LoadOrgMembers(ctx, userIDs)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] seeded environment created")
	err = bar.Finish()
	cobra.CheckErr(err)

	return getAllOrgMemberData(ctx, c, orgID)
}

// getAllOrgMemberData gets all the data from the seeded environment in a table format
func getAllOrgMemberData(ctx context.Context, c *seed.Client, orgID string) error {
	members, err := c.GetOrgMembersByOrgID(ctx, &openlaneclient.OrgMembershipWhereInput{
		OrganizationID: &orgID,
	})
	cobra.CheckErr(err)

	header := table.Row{"ID", "Email", "Role"}
	rows := []table.Row{}

	for _, om := range members.OrgMemberships.Edges {
		rows = append(rows, []interface{}{om.Node.ID, om.Node.User.Email, om.Node.Role})
	}

	// add empty row for spacing
	fmt.Println()

	createTableOutput("OrgMembers", header, rows)

	return nil
}
