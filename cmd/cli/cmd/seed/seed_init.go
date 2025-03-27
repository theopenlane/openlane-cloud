package seed

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"
	"github.com/theopenlane/openlane-cloud/internal/seed"
)

var seedInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init a new seeded environment",
	Long: `
	The init command will create a new demo environment with random data.
	The user must provide a PAT ID to authorize the root organization.
	A new API token will be created for the root organization and used to create the rest of the data.
	`,
	RunE: func(command *cobra.Command, _ []string) error {
		return initSeedData(command.Context())
	},
}

func init() {
	seedCmd.AddCommand(seedInitCmd)

	seedInitCmd.Flags().StringP("directory", "d", "demodata", "directory to save generated data")
	seedInitCmd.Flags().StringP("patid", "t", "", "personal access token ID to authorize the new root organization")
}

func initSeedData(ctx context.Context) error {
	c, err := newSeedClient()
	cobra.CheckErr(err)

	bar := progressbar.NewOptions(100, //nolint:mnd
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15), //nolint:mnd
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetDescription("[light_green]>[reset] creating seeded environment..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[light_green]=[reset]",
			SaucerHead:    "[light_green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	defer bar.Exit() //nolint:errcheck

	bar.Describe("[light_green]>[reset] registering users...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	userIDs, err := c.RegisterUsers(ctx)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] adding org members...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	err = c.LoadOrgMembers(ctx, userIDs)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] creating groups...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	err = c.LoadGroups(ctx)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] creating invites...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	err = c.LoadInvites(ctx)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] creating subscribers...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	err = c.LoadSubscribers(ctx)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] creating templates...")
	cmd.BarAdd(bar, 10) //nolint:mnd

	err = c.LoadTemplates(ctx)
	cobra.CheckErr(err)

	bar.Describe("[light_green]>[reset] seeded environment created")
	err = bar.Finish()
	cobra.CheckErr(err)

	return getAllData(ctx, c)
}

// createTableOutput creates a table output for the given data
func createTableOutput(title string, header table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetTitle(title)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(rows)
	t.Render()
}

// getAllData gets all the data from the seeded environment in a table format
func getAllData(ctx context.Context, c *seed.Client) error {
	// get all the orgs
	orgs, err := c.GetAllOrganizations(ctx)
	cobra.CheckErr(err)

	fmt.Println()
	fmt.Println("Seeded Environment Created:")

	header := table.Row{"ID", "Name", "Description", "PersonalOrg", "Children", "Members"}
	rows := []table.Row{}

	for _, org := range orgs.Organizations.Edges {
		rows = append(rows, []interface{}{org.Node.ID, org.Node.DisplayName, *org.Node.Description, *org.Node.PersonalOrg, len(org.Node.Children.Edges), len(org.Node.Members)})
	}

	createTableOutput("Organization", header, rows)

	groups, err := c.GetAllGroups(ctx)
	cobra.CheckErr(err)

	header = table.Row{"ID", "Name", "Description", "Visibility", "Members"}
	rows = []table.Row{}

	for _, group := range groups.Groups.Edges {
		rows = append(rows, []interface{}{group.Node.ID, group.Node.Name, *group.Node.Description, group.Node.Setting.Visibility, len(group.Node.Members)})
	}

	createTableOutput("Groups", header, rows)

	invites, err := c.GetAllInvites(ctx)
	cobra.CheckErr(err)

	header = table.Row{"ID", "Recipient", "Role", "Status"}
	rows = []table.Row{}

	for _, invite := range invites.Invites.Edges {
		rows = append(rows, []interface{}{invite.Node.ID, invite.Node.Recipient, invite.Node.Role, invite.Node.Status})
	}

	createTableOutput("Invites", header, rows)

	subscribers, err := c.GetAllSubscribers(ctx)
	cobra.CheckErr(err)

	header = table.Row{"ID", "Email", "Active", "Verified"}
	rows = []table.Row{}

	for _, sub := range subscribers.Subscribers.Edges {
		rows = append(rows, []interface{}{sub.Node.ID, sub.Node.Email, sub.Node.Active, sub.Node.VerifiedEmail})
	}

	createTableOutput("Subscribers", header, rows)

	templates, err := c.GetAllTemplates(ctx)
	cobra.CheckErr(err)

	header = table.Row{"ID", "Name"}
	rows = []table.Row{}

	for _, template := range templates.Templates.Edges {
		rows = append(rows, []interface{}{template.Node.ID, template.Node.Name})
	}

	createTableOutput("Templates", header, rows)

	return nil
}

// newSeedClient creates a new seed client, requiring a token to be set
func newSeedClient() (*seed.Client, error) {
	conf, err := seed.NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	if cmd.Config.String("token") == "" {
		return nil, cmd.ErrOpenlaneAPITokenMissing
	}

	if cmd.Config.String("directory") != "" {
		conf.Directory = cmd.Config.String("directory")
	}

	if cmd.Config.String("openlanehost") != "" {
		conf.OpenlaneHost = cmd.Config.String("openlanehost")
	}

	conf.Token = cmd.Config.String("token")

	return conf.NewClient()
}
