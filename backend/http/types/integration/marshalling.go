package integration

import (
	integrationtypes "ems.dev/backend/services/integration/types"
)

func MarshalIntegrationConfig(integration *integrationtypes.IntegrationConfig) *GetOrganizationIntegrationConfigsRequest {
	return &GetOrganizationIntegrationConfigsRequest{
		ID:             integration.ID,
		OrganizationID: integration.OrganizationID,
		ProviderName:   integration.ProviderName,
		ProviderType:   integration.ProviderType,
		Metadata:       integration.Metadata,
		LastSyncedAt:   integration.LastSyncedAt,
		CreatedAt:      integration.CreatedAt,
		UpdatedAt:      integration.UpdatedAt,
	}
}
