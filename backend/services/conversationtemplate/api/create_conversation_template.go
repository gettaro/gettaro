package api

import (
	"ems.dev/backend/services/conversationtemplate/types"
)

// CreateConversationTemplate creates a new conversation template
func (a *ConversationTemplateAPI) CreateConversationTemplate(params types.CreateConversationTemplateParams) (*types.ConversationTemplate, error) {
	return a.db.CreateConversationTemplate(params)
}
