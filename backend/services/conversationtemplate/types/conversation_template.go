package types

import (
	"time"

	"github.com/google/uuid"
)

// TemplateField represents a field in a conversation template
type TemplateField struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Type        string   `json:"type"` // text, textarea, select, checkbox, rating, etc.
	Required    bool     `json:"required"`
	Options     []string `json:"options,omitempty"` // For select/checkbox fields
	Placeholder string   `json:"placeholder,omitempty"`
	Order       int      `json:"order"`
}

// ConversationTemplate represents a conversation template
type ConversationTemplate struct {
	ID             uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID       `json:"organization_id" gorm:"type:uuid;not null"`
	Name           string          `json:"name" gorm:"not null"`
	Description    *string         `json:"description"`
	TemplateFields []TemplateField `json:"template_fields" gorm:"type:jsonb;not null;default:'[]'::jsonb"`
	IsActive       bool            `json:"is_active" gorm:"not null;default:true"`
	CreatedAt      time.Time       `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt      time.Time       `json:"updated_at" gorm:"not null;default:now()"`
}

// CreateConversationTemplateParams represents parameters for creating a conversation template
type CreateConversationTemplateParams struct {
	OrganizationID uuid.UUID       `json:"organization_id"`
	Name           string          `json:"name"`
	Description    *string         `json:"description"`
	TemplateFields []TemplateField `json:"template_fields"`
	IsActive       *bool           `json:"is_active"`
}

// UpdateConversationTemplateParams represents parameters for updating a conversation template
type UpdateConversationTemplateParams struct {
	ID             uuid.UUID        `json:"id"`
	Name           *string          `json:"name"`
	Description    *string          `json:"description"`
	TemplateFields *[]TemplateField `json:"template_fields"`
	IsActive       *bool            `json:"is_active"`
}

// ConversationTemplateSearchParams represents search parameters for conversation templates
type ConversationTemplateSearchParams struct {
	OrganizationID *uuid.UUID `json:"organization_id"`
	IsActive       *bool      `json:"is_active"`
}
