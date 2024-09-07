package prompts

import (
	"github.com/manifoldco/promptui"

	"github.com/theopenlane/openlane-cloud/cmd/cli/cmd"
)

func Name() (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return cmd.NewRequiredFieldMissingError("name")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:     "Name:",
		Templates: templates,
		Validate:  validate,
	}

	return prompt.Run()
}
