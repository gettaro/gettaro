package conversation

import (
	"ems.dev/backend/services/conversation/types"
	"gorm.io/datatypes"
)

// ListConversationsResponse represents the response for listing conversations
type ListConversationsResponse struct {
	Conversations []*types.Conversation `json:"conversations"`
}

// GetConversationResponse represents the response for getting a conversation
type GetConversationResponse struct {
	Conversation *types.Conversation `json:"conversation"`
}

// GetConversationWithDetailsResponse represents the response for getting a conversation with details
type GetConversationWithDetailsResponse struct {
	Conversation *types.ConversationWithDetails `json:"conversation"`
}

// CreateConversationRequest represents the request to create a conversation
type CreateConversationRequest struct {
	TemplateID       *string        `json:"template_id,omitempty"`
	Title            string         `json:"title" binding:"required"`
	DirectMemberID   string         `json:"direct_member_id" binding:"required"`
	ConversationDate *string        `json:"conversation_date,omitempty"` // ISO date string
	Content          datatypes.JSON `json:"content,omitempty"`
}

// UpdateConversationRequest represents the request to update a conversation
type UpdateConversationRequest struct {
	ConversationDate *string        `json:"conversation_date,omitempty"` // ISO date string
	Status           *string        `json:"status,omitempty"`
	Content          datatypes.JSON `json:"content,omitempty"`
}

// ListConversationsQuery represents query parameters for listing conversations
type ListConversationsQuery struct {
	ManagerMemberID *string `form:"manager_member_id"`
	DirectMemberID  *string `form:"direct_member_id"`
	TemplateID      *string `form:"template_id"`
	Status          *string `form:"status"`
	Limit           *int    `form:"limit"`
	Offset          *int    `form:"offset"`
}

// ConversationStatsResponse represents the response for conversation statistics
type ConversationStatsResponse struct {
	Stats map[string]int `json:"stats"`
}
