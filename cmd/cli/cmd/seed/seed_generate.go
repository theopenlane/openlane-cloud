package seed

import (
	"github.com/spf13/cobra"

	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"
	"github.com/theopenlane/openlane-cloud/internal/seed"
)

var (
	defaultObjectCount     = 10
	defaultInviteCount     = 5
	defaultSubscriberCount = 30
)

var seedGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate random data for seeded environment with a single root organization",
	RunE: func(_ *cobra.Command, _ []string) error {
		return generate()
	},
}

func init() {
	seedCmd.AddCommand(seedGenerateCmd)

	seedGenerateCmd.Flags().StringP("directory", "d", "demodata", "directory to save generated data")
	seedGenerateCmd.Flags().Int("users", defaultObjectCount, "number of users to generate")
	seedGenerateCmd.Flags().Int("groups", defaultObjectCount, "approximate number of groups to generate")
	seedGenerateCmd.Flags().Int("invites", defaultInviteCount, "number of invites to generate")
	seedGenerateCmd.Flags().Int("subscribers", defaultSubscriberCount, "number of subscribers to generate")
}

func generate() error {
	config, err := seed.NewDefaultConfig()
	cobra.CheckErr(err)

	if cmd.Config.String("directory") != "" {
		config.Directory = cmd.Config.String("directory")
	}

	config.NumUsers = cmd.Config.Int("users")
	config.NumGroups = cmd.Config.Int("groups")
	config.NumInvites = cmd.Config.Int("invites")
	config.NumSubscribers = cmd.Config.Int("subscribers")

	return config.GenerateData()
}
