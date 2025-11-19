package api

import (
	"context"

	"ems.dev/backend/services/conversation/database"
	"ems.dev/backend/services/conversation/types"
	templateapi "ems.dev/backend/services/conversationtemplate/api"
)

// ConversationAPIInterface defines the interface for conversation operations
type ConversationAPIInterface interface {
	CreateConversation(ctx context.Context, organizationID string, managerMemberID string, req *types.CreateConversationRequest) (*types.Conversation, error)
	GetConversation(ctx context.Context, id string) (*types.Conversation, error)
	GetConversationWithDetails(ctx context.Context, id string) (*types.ConversationWithDetails, error)
	ListConversations(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.Conversation, error)
	ListConversationsWithDetails(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.ConversationWithDetails, error)
	UpdateConversation(ctx context.Context, id string, req *types.UpdateConversationRequest) error
	DeleteConversation(ctx context.Context, id string) error
	GetConversationStats(ctx context.Context, organizationID string, managerMemberID *string) (map[string]int, error)
}

// ConversationAPI implements the conversation API
type ConversationAPI struct {
	db          database.DB
	templateApi templateapi.ConversationTemplateAPIInterface
}

// NewConversationAPI creates a new conversation API instance
func NewConversationAPI(db database.DB, templateApi templateapi.ConversationTemplateAPIInterface) *ConversationAPI {
	return &ConversationAPI{db: db, templateApi: templateApi}
}
