package api

import (
	"context"
)

// DeleteConversation deletes a conversation
func (a *ConversationAPI) DeleteConversation(ctx context.Context, id string) error {
	return a.db.DeleteConversation(ctx, id)
}
