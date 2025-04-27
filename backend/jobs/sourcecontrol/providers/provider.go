package providers

import (
	"context"

	"ems.dev/backend/services/integration/types"
)

// SourceControlProvider defines the interface for source control providers
type SourceControlProvider interface {
	// Name returns the unique identifier for this source control provider
	Name() string
	// SyncRepositories fetches and syncs data for the given repositories
	SyncRepositories(ctx context.Context, config *types.IntegrationConfig, repositories []string) error
}
