package openlane

import (
	"fmt"
)

var (
	// ErrAPITokenMissing is returned when the openlane API token is missing
	ErrAPITokenMissing = fmt.Errorf("token is required but not provided")
)
