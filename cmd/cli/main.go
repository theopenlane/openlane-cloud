package main

import (
	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"

	_ "github.com/theopenlane/openlane-cloud/cmd/cli/cmd/organization"
	_ "github.com/theopenlane/openlane-cloud/cmd/cli/cmd/seed"
)

func main() {
	cmd.Execute()
}
