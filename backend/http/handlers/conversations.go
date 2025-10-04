package handlers

import (
	"net/http"
	"time"

	"ems.dev/backend/http/types/conversation"
	"ems.dev/backend/http/utils"
	"ems.dev/backend/services/conversation/api"
	"ems.dev/backend/services/conversation/types"
	memberapi "ems.dev/backend/services/member/api"
	membertypes "ems.dev/backend/services/member/types"
	orgapi "ems.dev/backend/services/organization/api"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

type ConversationHandler struct {
	conversationApi api.ConversationAPIInterface
	memberApi       memberapi.MemberAPI
	orgApi          orgapi.OrganizationAPI
	userApi         userapi.UserAPI
}

func NewConversationHandler(conversationApi api.ConversationAPIInterface, memberApi memberapi.MemberAPI, orgApi orgapi.OrganizationAPI, userApi userapi.UserAPI) *ConversationHandler {
	return &ConversationHandler{
		conversationApi: conversationApi,
		memberApi:       memberApi,
		orgApi:          orgApi,
		userApi:         userApi,
	}
}

// RegisterRoutes registers all conversation routes
func (h *ConversationHandler) RegisterRoutes(router *gin.RouterGroup) {
	conversations := router.Group("/organizations/:id/conversations")
	{
		conversations.GET("", h.ListConversations)
		conversations.POST("", h.CreateConversation)
		conversations.GET("/stats", h.GetConversationStats)
	}

	conversation := router.Group("/conversations")
	{
		conversation.GET("/:id", h.GetConversation)
		conversation.GET("/:id/details", h.GetConversationWithDetails)
		conversation.PUT("/:id", h.UpdateConversation)
		conversation.DELETE("/:id", h.DeleteConversation)
	}
}

// ListConversations handles listing conversations for an organization
func (h *ConversationHandler) ListConversations(c *gin.Context) {
	organizationID := c.Param("id")

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &organizationID) {
		return
	}

	// Parse query parameters
	var query conversation.ListConversationsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to service types
	serviceQuery := &types.ListConversationsQuery{
		ManagerMemberID: query.ManagerMemberID,
		DirectMemberID:  query.DirectMemberID,
		TemplateID:      query.TemplateID,
		Status:          query.Status,
		Limit:           query.Limit,
		Offset:          query.Offset,
	}

	conversations, err := h.conversationApi.ListConversations(c.Request.Context(), organizationID, serviceQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation.ListConversationsResponse{
		Conversations: conversations,
	})
}

// GetConversation handles getting a single conversation
func (h *ConversationHandler) GetConversation(c *gin.Context) {
	conversationID := c.Param("id")

	// Get conversation first to check organization membership
	conv, err := h.conversationApi.GetConversation(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &conv.OrganizationID) {
		return
	}

	c.JSON(http.StatusOK, conversation.GetConversationResponse{
		Conversation: conv,
	})
}

// GetConversationWithDetails handles getting a conversation with related data
func (h *ConversationHandler) GetConversationWithDetails(c *gin.Context) {
	conversationID := c.Param("id")

	// Get conversation first to check organization membership
	conv, err := h.conversationApi.GetConversation(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &conv.OrganizationID) {
		return
	}

	// Get conversation with details
	convWithDetails, err := h.conversationApi.GetConversationWithDetails(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation.GetConversationWithDetailsResponse{
		Conversation: convWithDetails,
	})
}

// CreateConversation handles creating a new conversation
func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	organizationID := c.Param("id")

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &organizationID) {
		return
	}

	// Get current user
	ctxUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	user := ctxUser.(*usertypes.User)

	// Get user's member record for this organization
	members, err := h.memberApi.GetOrganizationMembers(c.Request.Context(), organizationID, &membertypes.OrganizationMemberParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userMember *membertypes.OrganizationMember
	for _, member := range members {
		if member.UserID == user.ID {
			userMember = &member
			break
		}
	}

	if userMember == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this organization"})
		return
	}

	var req conversation.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse conversation date if provided
	var conversationDate *time.Time
	if req.ConversationDate != nil && *req.ConversationDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.ConversationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation_date format, expected YYYY-MM-DD"})
			return
		}
		conversationDate = &parsedDate
	}

	// Convert to service types
	serviceReq := &types.CreateConversationRequest{
		TemplateID:       req.TemplateID,
		Title:            req.Title,
		DirectMemberID:   req.DirectMemberID,
		ConversationDate: conversationDate,
		Content:          req.Content,
	}

	conv, err := h.conversationApi.CreateConversation(c.Request.Context(), organizationID, userMember.ID, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conversation.GetConversationResponse{
		Conversation: conv,
	})
}

// UpdateConversation handles updating a conversation
func (h *ConversationHandler) UpdateConversation(c *gin.Context) {
	conversationID := c.Param("id")

	// Get conversation first to check ownership
	conv, err := h.conversationApi.GetConversation(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &conv.OrganizationID) {
		return
	}

	// Get current user
	ctxUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	user := ctxUser.(*usertypes.User)

	// Get user's member record for this organization
	members, err := h.memberApi.GetOrganizationMembers(c.Request.Context(), conv.OrganizationID, &membertypes.OrganizationMemberParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userMember *membertypes.OrganizationMember
	for _, member := range members {
		if member.UserID == user.ID {
			userMember = &member
			break
		}
	}

	if userMember == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this organization"})
		return
	}

	// Check if user is the manager of this conversation
	if userMember.ID != conv.ManagerMemberID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the conversation manager can update this conversation"})
		return
	}

	var req conversation.UpdateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse conversation date if provided
	var conversationDate *time.Time
	if req.ConversationDate != nil && *req.ConversationDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.ConversationDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversation_date format, expected YYYY-MM-DD"})
			return
		}
		conversationDate = &parsedDate
	}

	// Convert status if provided
	var status *types.ConversationStatus
	if req.Status != nil {
		convStatus := types.ConversationStatus(*req.Status)
		status = &convStatus
	}

	// Convert to service types
	serviceReq := &types.UpdateConversationRequest{
		ConversationDate: conversationDate,
		Status:           status,
		Content:          req.Content,
	}

	err = h.conversationApi.UpdateConversation(c.Request.Context(), conversationID, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "conversation updated successfully"})
}

// DeleteConversation handles deleting a conversation
func (h *ConversationHandler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")

	// Get conversation first to check ownership
	conv, err := h.conversationApi.GetConversation(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &conv.OrganizationID) {
		return
	}

	// Get current user
	ctxUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}
	user := ctxUser.(*usertypes.User)

	// Get user's member record for this organization
	members, err := h.memberApi.GetOrganizationMembers(c.Request.Context(), conv.OrganizationID, &membertypes.OrganizationMemberParams{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userMember *membertypes.OrganizationMember
	for _, member := range members {
		if member.UserID == user.ID {
			userMember = &member
			break
		}
	}

	if userMember == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this organization"})
		return
	}

	// Check if user is the manager of this conversation
	if userMember.ID != conv.ManagerMemberID {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the conversation manager can delete this conversation"})
		return
	}

	err = h.conversationApi.DeleteConversation(c.Request.Context(), conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "conversation deleted successfully"})
}

// GetConversationStats handles getting conversation statistics
func (h *ConversationHandler) GetConversationStats(c *gin.Context) {
	organizationID := c.Param("id")
	managerMemberID := c.Query("manager_member_id")

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &organizationID) {
		return
	}

	stats, err := h.conversationApi.GetConversationStats(c.Request.Context(), organizationID, &managerMemberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversation.ConversationStatsResponse{
		Stats: stats,
	})
}
