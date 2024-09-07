package organization

import (
	"github.com/spf13/cobra"

	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"
)

// organizationCmd represents the base organization command when called without any subcommands
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "the subcommands for working with the openlane organizations",
}

func init() {
	cmd.RootCmd.AddCommand(organizationCmd)
}
