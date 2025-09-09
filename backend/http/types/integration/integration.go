package integration

import (
	"time"

	integrationtypes "ems.dev/backend/services/integration/types"
	"gorm.io/datatypes"
)

type GetOrganizationIntegrationConfigsRequest struct {
	ID             string                                   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID string                                   `json:"organization_id"`
	ProviderName   integrationtypes.IntegrationProvider     `json:"provider_name"`
	ProviderType   integrationtypes.IntegrationProviderType `json:"provider_type"`
	Metadata       datatypes.JSON                           `json:"metadata"`
	LastSyncedAt   *time.Time                               `json:"last_synced_at"`
	CreatedAt      time.Time                                `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time                                `json:"updated_at"`
}
