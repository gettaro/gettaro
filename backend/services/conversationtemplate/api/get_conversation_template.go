package api

import (
	"ems.dev/backend/services/conversationtemplate/types"
	"github.com/google/uuid"
)

// GetConversationTemplate retrieves a conversation template by ID
func (a *ConversationTemplateAPI) GetConversationTemplate(id uuid.UUID) (*types.ConversationTemplate, error) {
	return a.db.GetConversationTemplate(id)
}
