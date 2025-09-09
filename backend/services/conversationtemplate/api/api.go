package api

import (
	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
)

// ConversationTemplateAPIInterface defines the interface for conversation template operations
type ConversationTemplateAPIInterface interface {
	CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error)
	GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error)
	ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error)
	UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error)
	DeleteConversationTemplate(id uuid.UUID) error
}
