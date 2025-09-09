package conversationtemplate

import (
	"github.com/google/uuid"
)

// TemplateField represents a field in a conversation template
type TemplateField struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Options     []string `json:"options,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	Order       int      `json:"order"`
}

// ConversationTemplateResponse represents a conversation template in HTTP responses
type ConversationTemplateResponse struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	Name           string          `json:"name"`
	Description    *string         `json:"description"`
	TemplateFields []TemplateField `json:"template_fields"`
	IsActive       bool            `json:"is_active"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

// CreateConversationTemplateRequest represents a request to create a conversation template
type CreateConversationTemplateRequest struct {
	Name           string          `json:"name" binding:"required"`
	Description    *string         `json:"description"`
	TemplateFields []TemplateField `json:"template_fields" binding:"required"`
	IsActive       *bool           `json:"is_active"`
}

// UpdateConversationTemplateRequest represents a request to update a conversation template
type UpdateConversationTemplateRequest struct {
	Name           *string          `json:"name"`
	Description    *string          `json:"description"`
	TemplateFields *[]TemplateField `json:"template_fields"`
	IsActive       *bool            `json:"is_active"`
}

// ListConversationTemplatesQuery represents query parameters for listing conversation templates
type ListConversationTemplatesQuery struct {
	IsActive *bool `form:"is_active"`
}

// ListConversationTemplatesResponse represents the response for listing conversation templates
type ListConversationTemplatesResponse struct {
	ConversationTemplates []ConversationTemplateResponse `json:"conversation_templates"`
}

// GetConversationTemplateResponse represents the response for getting a single conversation template
type GetConversationTemplateResponse struct {
	ConversationTemplate ConversationTemplateResponse `json:"conversation_template"`
}

// CreateConversationTemplateResponse represents the response for creating a conversation template
type CreateConversationTemplateResponse struct {
	ConversationTemplate ConversationTemplateResponse `json:"conversation_template"`
}

// UpdateConversationTemplateResponse represents the response for updating a conversation template
type UpdateConversationTemplateResponse struct {
	ConversationTemplate ConversationTemplateResponse `json:"conversation_template"`
}
