package types

import (
	"time"

	"gorm.io/datatypes"
)

// ConversationStatus represents the status of a conversation
type ConversationStatus string

const (
	ConversationStatusDraft     ConversationStatus = "draft"
	ConversationStatusCompleted ConversationStatus = "completed"
)

// Conversation represents a conversation between a manager and direct report
type Conversation struct {
	ID               string             `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrganizationID   string             `json:"organization_id"`
	TemplateID       *string            `json:"template_id,omitempty"`
	Title            string             `json:"title"`
	ManagerMemberID  string             `json:"manager_member_id"`
	DirectMemberID   string             `json:"direct_member_id"`
	ConversationDate *time.Time         `json:"conversation_date,omitempty"`
	Status           ConversationStatus `json:"status"`
	Content          datatypes.JSON     `json:"content,omitempty"` // Filled template data + template field definitions
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// ConversationWithDetails includes related data for display
type ConversationWithDetails struct {
	Conversation
	Template     *ConversationTemplate `json:"template,omitempty"`
	Manager      *OrganizationMember   `json:"manager,omitempty"`
	DirectReport *OrganizationMember   `json:"direct_report,omitempty"`
}

// CreateConversationRequest represents the request to create a conversation
type CreateConversationRequest struct {
	TemplateID       *string        `json:"template_id,omitempty"`
	Title            string         `json:"title" binding:"required"`
	DirectMemberID   string         `json:"direct_member_id" binding:"required"`
	ConversationDate *time.Time     `json:"conversation_date,omitempty"`
	Content          datatypes.JSON `json:"content,omitempty"` // Filled template data
}

// UpdateConversationRequest represents the request to update a conversation
type UpdateConversationRequest struct {
	ConversationDate *time.Time          `json:"conversation_date,omitempty"`
	Status           *ConversationStatus `json:"status,omitempty"`
	Content          datatypes.JSON      `json:"content,omitempty"` // Filled template data
}

// ListConversationsQuery represents query parameters for listing conversations
type ListConversationsQuery struct {
	ManagerMemberID *string `json:"manager_member_id,omitempty"`
	DirectMemberID  *string `json:"direct_member_id,omitempty"`
	TemplateID      *string `json:"template_id,omitempty"`
	Status          *string `json:"status,omitempty"`
	Limit           *int    `json:"limit,omitempty"`
	Offset          *int    `json:"offset,omitempty"`
}

// ConversationTemplate represents a conversation template (imported from conversationtemplate service)
type ConversationTemplate struct {
	ID             string          `json:"id"`
	OrganizationID string          `json:"organization_id"`
	Name           string          `json:"name"`
	Description    *string         `json:"description,omitempty"`
	TemplateFields []TemplateField `json:"template_fields"`
	IsActive       bool            `json:"is_active"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// TemplateField represents a field in a conversation template
type TemplateField struct {
	ID          string  `json:"id"`
	Label       string  `json:"label"`
	Type        string  `json:"type"`
	Required    bool    `json:"required"`
	Placeholder *string `json:"placeholder,omitempty"`
	Order       int     `json:"order"`
}

// OrganizationMember represents a member (imported from member service)
type OrganizationMember struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	OrganizationID string    `json:"organization_id"`
	IsOwner        bool      `json:"is_owner"`
	TitleID        *string   `json:"title_id,omitempty"`
	ManagerID      *string   `json:"manager_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
