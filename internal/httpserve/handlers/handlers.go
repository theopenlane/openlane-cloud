package handlers

import (
	"github.com/theopenlane/core/pkg/openlaneclient"
)

// Handler contains configuration options for handlers
type Handler struct {
	// IsTest is a flag to determine if the application is running in test mode and will mock external calls
	IsTest bool
	// ReadyChecks is a set of checkFuncs to determine if the application is "ready" upon startup
	ReadyChecks Checks
	// OpenlaneClient is the client to interact with the openlane API
	OpenlaneClient *openlaneclient.OpenlaneClient
}
