package providers

import (
	"context"

	"ems.dev/backend/services/integration/types"
)

// AICodeAssistantProvider defines the interface for AI code assistant providers
type AICodeAssistantProvider interface {
	// Name returns the unique identifier for this provider
	Name() string
	// SyncUsageData fetches and syncs usage data for the given integration config
	SyncUsageData(ctx context.Context, config *types.IntegrationConfig) error
}
