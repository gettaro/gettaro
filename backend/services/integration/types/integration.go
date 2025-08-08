package types

import (
	"time"

	"gorm.io/datatypes"
)

type IntegrationProvider string

const (
	IntegrationProviderGithub IntegrationProvider = "github"
)

type IntegrationProviderType string

const (
	IntegrationProviderTypeSourceControl     IntegrationProviderType = "SourceControl"
	IntegrationProviderTypeProjectManagement IntegrationProviderType = "ProjectManagement"
)

type IntegrationConfig struct {
	ID             string                  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID string                  `json:"organizationId"`
	ProviderName   IntegrationProvider     `json:"providerName"`
	ProviderType   IntegrationProviderType `json:"providerType"`
	EncryptedToken string                  `json:"encryptedToken"`
	Metadata       datatypes.JSON          `json:"metadata"`
	LastSyncedAt   *time.Time              `json:"lastSyncedAt"`
	CreatedAt      time.Time               `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time               `json:"updatedAt"`
}

type CreateIntegrationConfigRequest struct {
	ProviderName IntegrationProvider `json:"providerName" binding:"required"`
	Token        string              `json:"token" binding:"required"`
	Metadata     datatypes.JSON      `json:"metadata"`
}

type UpdateIntegrationConfigRequest struct {
	Token    string         `json:"token"`
	Metadata datatypes.JSON `json:"metadata"`
}
