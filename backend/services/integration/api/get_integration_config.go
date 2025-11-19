package api

import (
	"context"

	"ems.dev/backend/services/integration/types"
)

// GetIntegrationConfig retrieves an integration config by ID
func (a *Api) GetIntegrationConfig(ctx context.Context, id string) (*types.IntegrationConfig, error) {
	return a.db.GetIntegrationConfig(id)
}
