package integration

import (
	"time"

	integrationtypes "ems.dev/backend/services/integration/types"
	"gorm.io/datatypes"
)

type GetOrganizationIntegrationConfigsRequest struct {
	ID             string                                   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID string                                   `json:"organizationId"`
	ProviderName   integrationtypes.IntegrationProvider     `json:"providerName"`
	ProviderType   integrationtypes.IntegrationProviderType `json:"providerType"`
	Metadata       datatypes.JSON                           `json:"metadata"`
	LastSyncedAt   *time.Time                               `json:"lastSyncedAt"`
	CreatedAt      time.Time                                `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time                                `json:"updatedAt"`
}
