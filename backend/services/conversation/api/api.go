package api

import (
	"context"
	"encoding/json"

	"ems.dev/backend/services/conversation/database"
	"ems.dev/backend/services/conversation/types"
	templateapi "ems.dev/backend/services/conversationtemplate/api"
	"github.com/google/uuid"
	"gorm.io/datatypes"
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
	db          *database.ConversationDB
	templateApi templateapi.ConversationTemplateAPIInterface
}

// NewConversationAPI creates a new conversation API instance
func NewConversationAPI(db *database.ConversationDB, templateApi templateapi.ConversationTemplateAPIInterface) *ConversationAPI {
	return &ConversationAPI{db: db, templateApi: templateApi}
}

// CreateConversation creates a new conversation
func (a *ConversationAPI) CreateConversation(ctx context.Context, organizationID string, managerMemberID string, req *types.CreateConversationRequest) (*types.Conversation, error) {
	// Start with the provided content
	var contentJSON datatypes.JSON
	if req.Content != nil {
		contentJSON = req.Content
	} else {
		contentJSON = datatypes.JSON("{}")
	}

	// Parse content to add template metadata if needed
	var contentMap map[string]interface{}
	if err := json.Unmarshal(contentJSON, &contentMap); err != nil {
		contentMap = make(map[string]interface{})
	}

	// Determine the title - use provided title or default to template name
	title := req.Title
	if req.TemplateID != nil {
		templateID, err := uuid.Parse(*req.TemplateID)
		if err == nil {
			template, err := a.templateApi.GetConversationTemplate(templateID)
			if err == nil && template != nil {
				// Add template fields metadata
				contentMap["_template_fields"] = template.TemplateFields

				// If no title provided, use template name
				if title == "" {
					title = template.Name
				}
			}
		}
	}

	// Convert contentMap back to datatypes.JSON
	var err error
	contentJSON, err = json.Marshal(contentMap)
	if err != nil {
		return nil, err
	}

	conversation := &types.Conversation{
		OrganizationID:   organizationID,
		TemplateID:       req.TemplateID,
		Title:            title,
		ManagerMemberID:  managerMemberID,
		DirectMemberID:   req.DirectMemberID,
		ConversationDate: req.ConversationDate,
		Status:           types.ConversationStatusDraft,
		Content:          contentJSON,
	}

	err = a.db.CreateConversation(ctx, conversation)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

// GetConversation retrieves a conversation by ID
func (a *ConversationAPI) GetConversation(ctx context.Context, id string) (*types.Conversation, error) {
	return a.db.GetConversation(ctx, id)
}

// GetConversationWithDetails retrieves a conversation with related data
func (a *ConversationAPI) GetConversationWithDetails(ctx context.Context, id string) (*types.ConversationWithDetails, error) {
	return a.db.GetConversationWithDetails(ctx, id)
}

// ListConversations retrieves conversations with optional filters
func (a *ConversationAPI) ListConversations(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.Conversation, error) {
	return a.db.ListConversations(ctx, organizationID, query)
}

// ListConversationsWithDetails retrieves conversations with related data
func (a *ConversationAPI) ListConversationsWithDetails(ctx context.Context, organizationID string, query *types.ListConversationsQuery) ([]*types.ConversationWithDetails, error) {
	return a.db.ListConversationsWithDetails(ctx, organizationID, query)
}

// UpdateConversation updates a conversation
func (a *ConversationAPI) UpdateConversation(ctx context.Context, id string, req *types.UpdateConversationRequest) error {
	return a.db.UpdateConversation(ctx, id, req)
}

// DeleteConversation deletes a conversation
func (a *ConversationAPI) DeleteConversation(ctx context.Context, id string) error {
	return a.db.DeleteConversation(ctx, id)
}

// GetConversationStats returns statistics for conversations
func (a *ConversationAPI) GetConversationStats(ctx context.Context, organizationID string, managerMemberID *string) (map[string]int, error) {
	return a.db.GetConversationStats(ctx, organizationID, managerMemberID)
}
