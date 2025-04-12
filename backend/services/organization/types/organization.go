package types

import (
	"time"

	"gorm.io/datatypes"
)

type Organization struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string `gorm:"uniqueIndex"`
	Slug string `gorm:"uniqueIndex"`
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
