package api

import (
	"ems.dev/backend/services/conversationtemplate/types"
)

// ListConversationTemplates retrieves conversation templates based on search parameters
func (a *ConversationTemplateAPI) ListConversationTemplates(params types.ConversationTemplateSearchParams) ([]*types.ConversationTemplate, error) {
	return a.db.ListConversationTemplates(params)
}
