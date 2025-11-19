package api

import (
	"ems.dev/backend/services/conversationtemplate/database"
	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
)

// ConversationTemplateDBInterface defines the interface for database operations
type ConversationTemplateDBInterface interface {
	CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error)
	GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error)
	ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error)
	UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error)
	DeleteConversationTemplate(id uuid.UUID) error
}

// ConversationTemplateAPIInterface defines the interface for conversation template operations
type ConversationTemplateAPIInterface interface {
	CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error)
	GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error)
	ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error)
	UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error)
	DeleteConversationTemplate(id uuid.UUID) error
}

// ConversationTemplateAPI handles conversation template operations
type ConversationTemplateAPI struct {
	db ConversationTemplateDBInterface
}

// NewConversationTemplateAPI creates a new ConversationTemplateAPI instance
func NewConversationTemplateAPI(db *database.ConversationTemplateDatabase) *ConversationTemplateAPI {
	return &ConversationTemplateAPI{db: db}
}
