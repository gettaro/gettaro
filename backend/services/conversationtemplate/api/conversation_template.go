package api

import (
	"ems.dev/backend/services/conversationtemplate/database"
	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
)

// ConversationTemplateAPI handles conversation template operations
type ConversationTemplateAPI struct {
	db *database.ConversationTemplateDatabase
}

// NewConversationTemplateAPI creates a new ConversationTemplateAPI instance
func NewConversationTemplateAPI(db *database.ConversationTemplateDatabase) *ConversationTemplateAPI {
	return &ConversationTemplateAPI{db: db}
}

// CreateConversationTemplate creates a new conversation template
func (a *ConversationTemplateAPI) CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error) {
	return a.db.CreateConversationTemplate(params)
}

// GetConversationTemplate retrieves a conversation template by ID
func (a *ConversationTemplateAPI) GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error) {
	return a.db.GetConversationTemplate(id)
}

// ListConversationTemplates retrieves conversation templates based on search parameters
func (a *ConversationTemplateAPI) ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error) {
	return a.db.ListConversationTemplates(params)
}

// UpdateConversationTemplate updates an existing conversation template
func (a *ConversationTemplateAPI) UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error) {
	return a.db.UpdateConversationTemplate(params)
}

// DeleteConversationTemplate deletes a conversation template
func (a *ConversationTemplateAPI) DeleteConversationTemplate(id uuid.UUID) error {
	return a.db.DeleteConversationTemplate(id)
}
