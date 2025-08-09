package handlers

import (
	"net/http"

	httptypes "ems.dev/backend/http/types/integration"
	"ems.dev/backend/http/utils"
	intapi "ems.dev/backend/services/integration/api"
	inttypes "ems.dev/backend/services/integration/types"
	orgapi "ems.dev/backend/services/organization/api"
	"github.com/gin-gonic/gin"
)

type IntegrationHandler struct {
	integrationAPI intapi.IntegrationAPI
	orgAPI         orgapi.OrganizationAPI
}

func NewIntegrationHandler(integrationAPI intapi.IntegrationAPI, orgAPI orgapi.OrganizationAPI) *IntegrationHandler {
	return &IntegrationHandler{
		integrationAPI: integrationAPI,
		orgAPI:         orgAPI,
	}
}

// CreateIntegrationConfig handles the creation of a new integration config
func (h *IntegrationHandler) CreateIntegrationConfig(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgAPI, orgID) {
		return
	}

	var req inttypes.CreateIntegrationConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := h.integrationAPI.CreateIntegrationConfig(c.Request.Context(), orgID, &req)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"integration": config})
}

// GetIntegrationConfig handles retrieving a specific integration config
func (h *IntegrationHandler) GetIntegrationConfig(c *gin.Context) {
	id := c.Param("id")

	config, err := h.integrationAPI.GetIntegrationConfig(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &config.OrganizationID) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"integration": config})
}

// GetOrganizationIntegrationConfigs handles retrieving all integration configs for an organization
func (h *IntegrationHandler) GetOrganizationIntegrationConfigs(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgAPI, &orgID) {
		return
	}

	configs, err := h.integrationAPI.GetOrganizationIntegrationConfigs(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	marshalledConfigs := make([]*httptypes.GetOrganizationIntegrationConfigsRequest, len(configs))
	for i, config := range configs {
		marshalledConfigs[i] = httptypes.MarshalIntegrationConfig(&config)
	}

	c.JSON(http.StatusOK, gin.H{"integrations": marshalledConfigs})
}

// UpdateIntegrationConfig handles updating an existing integration config
func (h *IntegrationHandler) UpdateIntegrationConfig(c *gin.Context) {
	id := c.Param("id")

	config, err := h.integrationAPI.GetIntegrationConfig(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgAPI, config.OrganizationID) {
		return
	}

	var req inttypes.UpdateIntegrationConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedConfig, err := h.integrationAPI.UpdateIntegrationConfig(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"integration": updatedConfig})
}

// DeleteIntegrationConfig handles deleting an integration config
func (h *IntegrationHandler) DeleteIntegrationConfig(c *gin.Context) {
	if len(c.Params) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid number of parameters"})
		return
	}

	id := c.Params[1].Value

	config, err := h.integrationAPI.GetIntegrationConfig(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgAPI, config.OrganizationID) {
		return
	}

	if err := h.integrationAPI.DeleteIntegrationConfig(c.Request.Context(), id); err != nil {
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
