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

type IntegrationConfig struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrganizationID string
	ProviderName   string
	ProviderType   string
	EncryptedToken string
	Metadata       datatypes.JSON
	LastSyncedAt   *time.Time
	CreatedAt      time.Time `gorm:"default:now()"`
	UpdatedAt      time.Time

	// Unique constraint
	UniqueOrgProvider string `gorm:"uniqueIndex:idx_org_provider"`
}

type OrganizationMember struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"userId"`
	OrganizationID string    `json:"organizationId"`
	IsOwner        bool      `json:"isOwner"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
