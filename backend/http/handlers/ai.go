package handlers

import (
	"net/http"
	"strconv"

	"ems.dev/backend/http/types/ai"
	"ems.dev/backend/http/utils"
	apiai "ems.dev/backend/services/ai/api"
	aitypes "ems.dev/backend/services/ai/types"
	orgapi "ems.dev/backend/services/organization/api"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	aiService apiai.AIServiceInterface
	orgAPI    orgapi.OrganizationAPI
}

func NewAIHandler(aiService apiai.AIServiceInterface, orgAPI orgapi.OrganizationAPI) *AIHandler {
	return &AIHandler{
		aiService: aiService,
		orgAPI:    orgAPI,
	}
}

// RegisterRoutes registers AI routes
func (h *AIHandler) RegisterRoutes(rg *gin.RouterGroup) {
	organizations := rg.Group("/organizations/:id")
	{
		ai := organizations.Group("/ai")
		{
			ai.POST("/query", h.QueryAI)
			ai.GET("/history", h.GetQueryHistory)
			ai.GET("/stats", h.GetQueryStats)
		}
	}
}

// QueryAI handles AI query requests
func (h *AIHandler) QueryAI(c *gin.Context) {
	// Get organization ID from context
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check organization membership
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Parse request
	var req ai.AIQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Organization ID comes from the URL context

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to service request
	serviceReq := &aitypes.AIQueryRequest{
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		OrganizationID: orgID,
		Query:          req.Query,
		Context:        req.Context,
		AdditionalData: req.AdditionalData,
	}

	// Process query
	response, err := h.aiService.Query(c.Request.Context(), serviceReq, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to HTTP response
	httpResponse := ai.AIQueryResponse{
		Answer:      response.Answer,
		Sources:     response.Sources,
		Confidence:  response.Confidence,
		RelatedData: response.RelatedData,
		Suggestions: response.Suggestions,
	}

	c.JSON(http.StatusOK, gin.H{"ai_response": httpResponse})
}

// GetQueryHistory handles getting AI query history
func (h *AIHandler) GetQueryHistory(c *gin.Context) {
	// Get organization ID from context
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check organization membership
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	// Get query history
	userIDStr := userID.(string)
	history, err := h.aiService.GetQueryHistory(c.Request.Context(), orgID, &userIDStr, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to HTTP response
	httpHistory := make([]ai.AIQueryHistoryItem, len(history))
	for i, item := range history {
		httpHistory[i] = ai.AIQueryHistoryItem{
			ID:             item.ID,
			OrganizationID: item.OrganizationID,
			UserID:         item.UserID,
			EntityType:     item.EntityType,
			EntityID:       item.EntityID,
			Query:          item.Query,
			Answer:         item.Answer,
			Context:        item.Context,
			Confidence:     item.Confidence,
			Sources:        item.Sources,
			CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"queries": httpHistory})
}

// GetQueryStats handles getting AI query statistics
func (h *AIHandler) GetQueryStats(c *gin.Context) {
	// Get organization ID from context
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check organization membership
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Parse query parameters
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 30
	}

	// Get query stats
	userIDStr := userID.(string)
	stats, err := h.aiService.GetQueryStats(c.Request.Context(), orgID, &userIDStr, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to HTTP response
	httpStats := ai.AIQueryStatsResponse{
		TotalQueries:      stats.TotalQueries,
		QueriesByEntity:   stats.QueriesByEntity,
		QueriesByContext:  stats.QueriesByContext,
		AverageConfidence: stats.AverageConfidence,
		RecentQueries:     make([]ai.AIQueryHistoryItem, len(stats.RecentQueries)),
	}

	for i, item := range stats.RecentQueries {
		httpStats.RecentQueries[i] = ai.AIQueryHistoryItem{
			ID:             item.ID,
			OrganizationID: item.OrganizationID,
			UserID:         item.UserID,
			EntityType:     item.EntityType,
			EntityID:       item.EntityID,
			Query:          item.Query,
			Answer:         item.Answer,
			Context:        item.Context,
			Confidence:     item.Confidence,
			Sources:        item.Sources,
			CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	c.JSON(http.StatusOK, gin.H{"stats": httpStats})
}
