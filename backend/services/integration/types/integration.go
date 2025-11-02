package types

import (
	"time"

	"gorm.io/datatypes"
)

type IntegrationProvider string

const (
	IntegrationProviderGithub IntegrationProvider = "github"
	IntegrationProviderCursor IntegrationProvider = "cursor"
)

type IntegrationProviderType string

const (
	IntegrationProviderTypeSourceControl     IntegrationProviderType = "SourceControl"
	IntegrationProviderTypeProjectManagement IntegrationProviderType = "ProjectManagement"
	IntegrationProviderTypeAICodeAssistant    IntegrationProviderType = "AICodeAssistant"
)

type IntegrationConfig struct {
	ID             string                  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID string                  `json:"organization_id"`
	ProviderName   IntegrationProvider     `json:"provider_name"`
	ProviderType   IntegrationProviderType `json:"provider_type"`
	EncryptedToken string                  `json:"encrypted_token"`
	Metadata       datatypes.JSON          `json:"metadata"`
	LastSyncedAt   *time.Time              `json:"last_synced_at"`
	CreatedAt      time.Time               `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type CreateIntegrationConfigRequest struct {
	ProviderName IntegrationProvider `json:"provider_name" binding:"required"`
	Token        string              `json:"token" binding:"required"`
	Metadata     datatypes.JSON      `json:"metadata"`
}

type UpdateIntegrationConfigRequest struct {
	Token    string         `json:"token"`
	Metadata datatypes.JSON `json:"metadata"`
}
