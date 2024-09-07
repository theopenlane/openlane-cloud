package seed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/openlane-cloud/internal/seed"
)

func TestNewDefaultConfig(t *testing.T) {
	conf, err := seed.NewDefaultConfig()
	require.NoError(t, err)

	// Check default values
	assert.Equal(t, "demodata", conf.Directory)
	assert.Equal(t, "", conf.Token)
	assert.Equal(t, 1, conf.NumOrganizations)
	assert.Equal(t, 10, conf.NumUsers)
	assert.Equal(t, 10, conf.NumGroups)
	assert.Equal(t, 5, conf.NumInvites)
}
