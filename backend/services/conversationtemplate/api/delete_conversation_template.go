package api

import (
	"github.com/google/uuid"
)

// DeleteConversationTemplate deletes a conversation template
func (a *ConversationTemplateAPI) DeleteConversationTemplate(id uuid.UUID) error {
	return a.db.DeleteConversationTemplate(id)
}
