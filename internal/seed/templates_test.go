package seed

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemplates(t *testing.T) {
	// Call the function being tested
	templates, err := getTemplates(templateDirectory)
	require.NoError(t, err)

	// Check the number of templates
	expectedNumTemplates := 1
	assert.Len(t, templates, expectedNumTemplates)

	// Check the template name and JSON config
	for _, template := range templates {
		assert.NotEmpty(t, template.Name)
		assert.NotEmpty(t, template.JSONConfig)

		foundValidSchema := false

		for key, value := range template.JSONConfig {
			if key == "$schema" {
				foundValidSchema = true

				assert.Equal(t, "https://json-schema.org/draft/2020-12/schema", value)

				break
			}
		}

		assert.True(t, foundValidSchema)
	}

	// Call the function being tested, but include an invalid directory
	templates, err = getTemplates("invalid")
	require.Error(t, err)
	assert.Nil(t, templates)
}
