package handlers

import (
	"fmt"
	"net/http"

	"ems.dev/backend/http/utils"
	intapi "ems.dev/backend/services/integration/api"
	inttypes "ems.dev/backend/services/integration/types"
	orgapi "ems.dev/backend/services/organization/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

type IntegrationHandler struct {
	integrationAPI intapi.API
	orgAPI         orgapi.OrganizationAPI
}

func NewIntegrationHandler(integrationAPI intapi.API, orgAPI orgapi.OrganizationAPI) *IntegrationHandler {
	return &IntegrationHandler{
		integrationAPI: integrationAPI,
		orgAPI:         orgAPI,
	}
}

// getUserFromContext extracts the user from the request context and returns it
func (h *IntegrationHandler) getUserFromContext(c *gin.Context) (*usertypes.User, error) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return nil, fmt.Errorf("user not found in context: %w", err)
	}
	return user, nil
}

// getOrganizationIDFromContext extracts the organization ID from the request context and returns it
func (h *IntegrationHandler) getOrganizationIDFromContext(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("organization ID is required")
	}
	return id, nil
}

// CreateIntegrationConfig handles the creation of a new integration config
func (h *IntegrationHandler) CreateIntegrationConfig(c *gin.Context) {
	orgID, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgAPI, orgID) {
		return
	}

	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req inttypes.CreateIntegrationConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.integrationAPI.CreateIntegrationConfig(c.Request.Context(), orgID, user.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"integration": config})
}

// GetIntegrationConfig handles retrieving a specific integration config
func (h *IntegrationHandler) GetIntegrationConfig(c *gin.Context) {
	id := c.Param("id")
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	config, err := h.integrationAPI.GetIntegrationConfig(c.Request.Context(), id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationMembership(c, h.orgAPI, &config.OrganizationID) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"integration": config})
}

// GetOrganizationIntegrationConfigs handles retrieving all integration configs for an organization
func (h *IntegrationHandler) GetOrganizationIntegrationConfigs(c *gin.Context) {
	orgID, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	configs, err := h.integrationAPI.GetOrganizationIntegrationConfigs(c.Request.Context(), orgID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"integrations": configs})
}

// UpdateIntegrationConfig handles updating an existing integration config
func (h *IntegrationHandler) UpdateIntegrationConfig(c *gin.Context) {
	id := c.Param("id")
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	config, err := h.integrationAPI.GetIntegrationConfig(c.Request.Context(), id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgAPI, config.OrganizationID) {
		return
	}

	var req inttypes.UpdateIntegrationConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedConfig, err := h.integrationAPI.UpdateIntegrationConfig(c.Request.Context(), id, user.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"integration": updatedConfig})
}

// DeleteIntegrationConfig handles deleting an integration config
func (h *IntegrationHandler) DeleteIntegrationConfig(c *gin.Context) {
	id := c.Param("id")
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	config, err := h.integrationAPI.GetIntegrationConfig(c.Request.Context(), id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !utils.CheckOrganizationOwnership(c, h.orgAPI, config.OrganizationID) {
		return
	}

	if err := h.integrationAPI.DeleteIntegrationConfig(c.Request.Context(), id, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registers the integration routes
func (h *IntegrationHandler) RegisterRoutes(router *gin.RouterGroup) {
	integrations := router.Group("/organizations/:id/integrations")
	{
		integrations.POST("", h.CreateIntegrationConfig)
		integrations.GET("", h.GetOrganizationIntegrationConfigs)
		integrations.GET("/:id", h.GetIntegrationConfig)
		integrations.PUT("/:id", h.UpdateIntegrationConfig)
		integrations.DELETE("/:id", h.DeleteIntegrationConfig)
	}
}
