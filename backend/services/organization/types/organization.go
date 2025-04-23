package types

import (
	"time"

	"gorm.io/datatypes"
)

type Organization struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug" gorm:"uniqueIndex"`
	IsOwner   bool      `json:"isOwner" gorm:"-"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserOrganization struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	IsOwner        bool      `json:"isOwner" gorm:"-"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type IntegrationProvider string

const (
	IntegrationProviderGithub IntegrationProvider = "github"
)

type IntegrationConfig struct {
	ID             string              `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID string              `json:"organizationId"`
	ProviderName   IntegrationProvider `json:"providerName"`
	ProviderType   string              `json:"providerType"`
	EncryptedToken string              `json:"encryptedToken"`
	Metadata       datatypes.JSON      `json:"metadata"`
	LastSyncedAt   *time.Time          `json:"lastSyncedAt"`
	CreatedAt      time.Time           `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time           `json:"updatedAt"`
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

type OrganizationMember struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	IsOwner        bool      `json:"isOwner"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
