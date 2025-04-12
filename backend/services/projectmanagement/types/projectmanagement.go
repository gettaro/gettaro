package types

import (
	"time"

	"gorm.io/datatypes"
)

type ProjectManagementAccount struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string
	OrganizationID *string
	ProviderName   string
	ProviderID     string
	Metadata       datatypes.JSON
	LastSyncedAt   *time.Time
}

type PMTicket struct {
	ID                         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProjectManagementAccountID string
	ProviderID                 string
	Title                      string
	Epic                       *string
	Labels                     []string `gorm:"type:text[]"`
	Status                     string
	StoryPoints                *int
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	CompletedAt                *time.Time
}
