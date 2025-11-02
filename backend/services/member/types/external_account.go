package types

import (
	"time"

	"gorm.io/datatypes"
)

// ExternalAccountType represents the type of external account
type ExternalAccountType string

const (
	ExternalAccountTypeSourceControl      ExternalAccountType = "sourcecontrol"
	ExternalAccountTypeAICodeAssistant ExternalAccountType = "ai-code-assistant"
	// Future types can be added here:
	// ExternalAccountTypeJira ExternalAccountType = "jira"
	// ExternalAccountTypeSlack ExternalAccountType = "slack"
)

// ExternalAccount represents an external account linked to a member
type ExternalAccount struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	MemberID       *string        `gorm:"type:uuid" json:"member_id,omitempty"`
	OrganizationID *string        `gorm:"type:uuid" json:"organization_id,omitempty"`
	AccountType    string         `gorm:"type:varchar(50);check:account_type IN ('sourcecontrol','ai-code-assistant')" json:"account_type"`
	ProviderName   string         `gorm:"type:varchar(255)" json:"provider_name"`
	ProviderID     string         `gorm:"type:varchar(255)" json:"provider_id"`
	Username       string         `gorm:"type:varchar(255)" json:"username"`
	Metadata       datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	LastSyncedAt   *time.Time     `gorm:"type:timestamp with time zone" json:"last_synced_at,omitempty"`
	CreatedAt      time.Time      `gorm:"type:timestamp with time zone;default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamp with time zone;default:current_timestamp" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (ExternalAccount) TableName() string {
	return "member_external_accounts"
}

// ExternalAccountParams represents parameters for querying external accounts
type ExternalAccountParams struct {
	ExternalAccountIDs []string `json:"external_account_ids"`
	OrganizationID     string   `json:"organization_id"`
	MemberIDs          []string `json:"member_ids"`
	Usernames          []string `json:"usernames"`
	AccountType        *string  `json:"account_type,omitempty"` // Optional filter by account type
}

