package api

import (
	"context"

	"ems.dev/backend/services/integration/types"
)

// GetOrganizationIntegrationConfigs retrieves all integration configs for an organization
func (a *Api) GetOrganizationIntegrationConfigs(ctx context.Context, orgID string) ([]types.IntegrationConfig, error) {
	return a.db.GetOrganizationIntegrationConfigs(orgID)
}
