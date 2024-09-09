package seed

import (
	"github.com/spf13/cobra"

	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"
)

// seedCmd represents the base seed command when called without any subcommands
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "the subcommands for creating demo data in openlane",
}

func init() {
	cmd.RootCmd.AddCommand(seedCmd)
}
