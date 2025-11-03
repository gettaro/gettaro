package handlers

import (
	"net/http"
	"time"

	"ems.dev/backend/http/utils"
	aicodeassistantapi "ems.dev/backend/services/aicodeassistant/api"
	aicodeassistanttypes "ems.dev/backend/services/aicodeassistant/types"
	organizationapi "ems.dev/backend/services/organization/api"

	"github.com/gin-gonic/gin"
)

type AICodeAssistantHandler struct {
	aiCodeAssistantAPI aicodeassistantapi.AICodeAssistantAPI
	orgAPI             organizationapi.OrganizationAPI
}

func NewAICodeAssistantHandler(aiCodeAssistantAPI aicodeassistantapi.AICodeAssistantAPI, orgAPI organizationapi.OrganizationAPI) *AICodeAssistantHandler {
	return &AICodeAssistantHandler{
		aiCodeAssistantAPI: aiCodeAssistantAPI,
		orgAPI:              orgAPI,
	}
}

// GetOrganizationAICodeAssistantUsage handles retrieving AI code assistant usage metrics for an organization
// Query Parameters:
// - externalAccountIds: Optional list of external account IDs to filter by
// - toolName: Optional tool name to filter by (e.g., "cursor", "claude-code")
// - startDate: Optional start date in format "2006-01-02"
// - endDate: Optional end date in format "2006-01-02"
func (h *AICodeAssistantHandler) GetOrganizationAICodeAssistantUsage(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Parse query parameters
	params := &aicodeassistanttypes.AICodeAssistantDailyMetricParams{
		OrganizationID: orgID,
	}

	// Parse external account IDs if provided
	if externalAccountIds := c.QueryArray("externalAccountIds"); len(externalAccountIds) > 0 {
		params.ExternalAccountIDs = externalAccountIds
	}

	// Parse tool name if provided
	if toolName := c.Query("toolName"); toolName != "" {
		params.ToolName = &toolName
	}

	// Parse date range if provided
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format, expected YYYY-MM-DD"})
			return
		}
		params.StartDate = &startDate
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format, expected YYYY-MM-DD"})
			return
		}
		params.EndDate = &endDate
	}

	// Get daily metrics
	metrics, err := h.aiCodeAssistantAPI.GetDailyMetrics(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}

// GetOrganizationAICodeAssistantUsageStats handles retrieving aggregated statistics for an organization
// Query Parameters: Same as GetOrganizationAICodeAssistantUsage
func (h *AICodeAssistantHandler) GetOrganizationAICodeAssistantUsageStats(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Parse query parameters (same as GetOrganizationAICodeAssistantUsage)
	params := &aicodeassistanttypes.AICodeAssistantDailyMetricParams{
		OrganizationID: orgID,
	}

	if externalAccountIds := c.QueryArray("externalAccountIds"); len(externalAccountIds) > 0 {
		params.ExternalAccountIDs = externalAccountIds
	}

	if toolName := c.Query("toolName"); toolName != "" {
		params.ToolName = &toolName
	}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format, expected YYYY-MM-DD"})
			return
		}
		params.StartDate = &startDate
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format, expected YYYY-MM-DD"})
			return
		}
		params.EndDate = &endDate
	}

	// Get usage stats
	stats, err := h.aiCodeAssistantAPI.GetUsageStats(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetMemberAICodeAssistantUsage handles retrieving AI code assistant usage for a specific member
// Query Parameters:
// - toolName: Optional tool name to filter by
// - startDate: Optional start date in format "2006-01-02"
// - endDate: Optional end date in format "2006-01-02"
func (h *AICodeAssistantHandler) GetMemberAICodeAssistantUsage(c *gin.Context) {
	memberID := c.Param("memberId")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member ID is required"})
		return
	}

	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Parse query parameters
	params := &aicodeassistanttypes.AICodeAssistantMemberMetricsParams{}

	if toolName := c.Query("toolName"); toolName != "" {
		params.ToolName = &toolName
	}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format, expected YYYY-MM-DD"})
			return
		}
		params.StartDate = &startDate
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format, expected YYYY-MM-DD"})
			return
		}
		params.EndDate = &endDate
	}

	// Get daily metrics (service layer handles member -> external account resolution)
	metrics, err := h.aiCodeAssistantAPI.GetMemberDailyMetrics(c.Request.Context(), orgID, memberID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}

// GetMemberAICodeAssistantUsageStats handles retrieving aggregated statistics for a specific member
// Query Parameters:
// - toolName: Optional tool name to filter by
// - startDate: Optional start date in format "2006-01-02"
// - endDate: Optional end date in format "2006-01-02"
func (h *AICodeAssistantHandler) GetMemberAICodeAssistantUsageStats(c *gin.Context) {
	memberID := c.Param("memberId")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member ID is required"})
		return
	}

	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	// Parse query parameters
	params := &aicodeassistanttypes.AICodeAssistantMemberMetricsParams{}

	if toolName := c.Query("toolName"); toolName != "" {
		params.ToolName = &toolName
	}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid startDate format, expected YYYY-MM-DD"})
			return
		}
		params.StartDate = &startDate
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid endDate format, expected YYYY-MM-DD"})
			return
		}
		params.EndDate = &endDate
	}

	// Get usage stats (service layer handles member -> external account resolution)
	stats, err := h.aiCodeAssistantAPI.GetMemberUsageStats(c.Request.Context(), orgID, memberID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// RegisterRoutes registers the AI code assistant routes
func (h *AICodeAssistantHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Organization-level routes
	orgRoutes := router.Group("/organizations/:id")
	{
		orgRoutes.GET("/ai-code-assistant", h.GetOrganizationAICodeAssistantUsage)
		orgRoutes.GET("/ai-code-assistant/stats", h.GetOrganizationAICodeAssistantUsageStats)
	}

	// Member-level routes (under organization)
	memberRoutes := router.Group("/organizations/:id/members/:memberId")
	{
		memberRoutes.GET("/ai-code-assistant", h.GetMemberAICodeAssistantUsage)
		memberRoutes.GET("/ai-code-assistant/stats", h.GetMemberAICodeAssistantUsageStats)
	}
}
