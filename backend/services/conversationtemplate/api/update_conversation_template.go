package api

import (
	"ems.dev/backend/services/conversationtemplate/types"
)

// UpdateConversationTemplate updates an existing conversation template
func (a *ConversationTemplateAPI) UpdateConversationTemplate(params types.UpdateConversationTemplateParams) (*types.ConversationTemplate, error) {
	return a.db.UpdateConversationTemplate(params)
}
