package handlers

import (
	"net/http"
	"time"

	"ems.dev/backend/http/types/conversationtemplate"
	"ems.dev/backend/http/utils"
	"ems.dev/backend/services/conversationtemplate/api"
	"ems.dev/backend/services/conversationtemplate/types"
	orgapi "ems.dev/backend/services/organization/api"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConversationTemplateHandler handles HTTP requests for conversation templates
type ConversationTemplateHandler struct {
	conversationTemplateApi api.ConversationTemplateAPIInterface
	orgApi                  orgapi.OrganizationAPI
}

// NewConversationTemplateHandler creates a new ConversationTemplateHandler
func NewConversationTemplateHandler(conversationTemplateApi api.ConversationTemplateAPIInterface, orgApi orgapi.OrganizationAPI) *ConversationTemplateHandler {
	return &ConversationTemplateHandler{
		conversationTemplateApi: conversationTemplateApi,
		orgApi:                  orgApi,
	}
}

// RegisterRoutes registers all conversation template routes
func (h *ConversationTemplateHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Organization-specific conversation templates
	router.GET("/organizations/:id/conversation-templates", h.ListConversationTemplates)
	router.POST("/organizations/:id/conversation-templates", h.CreateConversationTemplate)

	// Individual conversation template operations
	router.GET("/conversation-templates/:id", h.GetConversationTemplate)
	router.PUT("/conversation-templates/:id", h.UpdateConversationTemplate)
	router.DELETE("/conversation-templates/:id", h.DeleteConversationTemplate)
}

// ListConversationTemplates handles GET /organizations/:id/conversation-templates
func (h *ConversationTemplateHandler) ListConversationTemplates(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	organizationID, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	// Parse query parameters
	var query conversationtemplate.ListConversationTemplatesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build search parameters
	searchParams := types.ConversationTemplateSearchParams{
		OrganizationID: &organizationID,
		IsActive:       query.IsActive,
	}

	// Get conversation templates
	templates, err := h.conversationTemplateApi.ListConversationTemplates(searchParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map to response format
	responseTemplates := make([]conversationtemplate.ConversationTemplateResponse, len(templates))
	for i, template := range templates {
		responseTemplates[i] = h.mapConversationTemplateToResponse(template)
	}

	response := conversationtemplate.ListConversationTemplatesResponse{
		ConversationTemplates: responseTemplates,
	}

	c.JSON(http.StatusOK, response)
}

// GetConversationTemplate handles GET /conversation-templates/:id
func (h *ConversationTemplateHandler) GetConversationTemplate(c *gin.Context) {
	// Get template ID from URL parameter
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	// Get conversation template
	template, err := h.conversationTemplateApi.GetConversationTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	orgIDStr := template.OrganizationID.String()
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgIDStr) {
		return
	}

	response := conversationtemplate.GetConversationTemplateResponse{
		ConversationTemplate: h.mapConversationTemplateToResponse(template),
	}

	c.JSON(http.StatusOK, response)
}

// CreateConversationTemplate handles POST /organizations/:id/conversation-templates
func (h *ConversationTemplateHandler) CreateConversationTemplate(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	organizationID, err := uuid.Parse(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	// Parse request body
	var req conversationtemplate.CreateConversationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map template fields
	templateFields := make([]types.TemplateField, len(req.TemplateFields))
	for i, field := range req.TemplateFields {
		templateFields[i] = types.TemplateField{
			ID:          field.ID,
			Label:       field.Label,
			Type:        field.Type,
			Required:    field.Required,
			Options:     field.Options,
			Placeholder: field.Placeholder,
			Order:       field.Order,
		}
	}

	// Create conversation template
	createParams := types.CreateConversationTemplateParams{
		OrganizationID: organizationID,
		Name:           req.Name,
		Description:    req.Description,
		TemplateFields: templateFields,
		IsActive:       req.IsActive,
	}

	template, err := h.conversationTemplateApi.CreateConversationTemplate(createParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := conversationtemplate.CreateConversationTemplateResponse{
		ConversationTemplate: h.mapConversationTemplateToResponse(template),
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateConversationTemplate handles PUT /conversation-templates/:id
func (h *ConversationTemplateHandler) UpdateConversationTemplate(c *gin.Context) {
	// Get template ID from URL parameter
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	// Get conversation template to check organization
	existingTemplate, err := h.conversationTemplateApi.GetConversationTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, existingTemplate.OrganizationID.String()) {
		return
	}

	// Parse request body
	var req conversationtemplate.UpdateConversationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Map template fields if provided
	var templateFields *[]types.TemplateField
	if req.TemplateFields != nil {
		fields := make([]types.TemplateField, len(*req.TemplateFields))
		for i, field := range *req.TemplateFields {
			fields[i] = types.TemplateField{
				ID:          field.ID,
				Label:       field.Label,
				Type:        field.Type,
				Required:    field.Required,
				Options:     field.Options,
				Placeholder: field.Placeholder,
				Order:       field.Order,
			}
		}
		templateFields = &fields
	}

	// Update conversation template
	updateParams := types.UpdateConversationTemplateParams{
		ID:             templateID,
		Name:           req.Name,
		Description:    req.Description,
		TemplateFields: templateFields,
		IsActive:       req.IsActive,
	}

	template, err := h.conversationTemplateApi.UpdateConversationTemplate(updateParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := conversationtemplate.UpdateConversationTemplateResponse{
		ConversationTemplate: h.mapConversationTemplateToResponse(template),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteConversationTemplate handles DELETE /conversation-templates/:id
func (h *ConversationTemplateHandler) DeleteConversationTemplate(c *gin.Context) {
	// Get template ID from URL parameter
	templateIDStr := c.Param("id")
	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid template ID"})
		return
	}

	// Get conversation template to check organization
	existingTemplate, err := h.conversationTemplateApi.GetConversationTemplate(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, existingTemplate.OrganizationID.String()) {
		return
	}

	// Delete conversation template
	if err := h.conversationTemplateApi.DeleteConversationTemplate(templateID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// mapConversationTemplateToResponse converts a types.ConversationTemplate to ConversationTemplateResponse
func (h *ConversationTemplateHandler) mapConversationTemplateToResponse(template *types.ConversationTemplate) conversationtemplate.ConversationTemplateResponse {
	// Map template fields
	responseFields := make([]conversationtemplate.TemplateField, len(template.TemplateFields))
	for i, field := range template.TemplateFields {
		responseFields[i] = conversationtemplate.TemplateField{
			ID:          field.ID,
			Label:       field.Label,
			Type:        field.Type,
			Required:    field.Required,
			Options:     field.Options,
			Placeholder: field.Placeholder,
			Order:       field.Order,
		}
	}

	return conversationtemplate.ConversationTemplateResponse{
		ID:             template.ID,
		OrganizationID: template.OrganizationID,
		Name:           template.Name,
		Description:    template.Description,
		TemplateFields: responseFields,
		IsActive:       template.IsActive,
		CreatedAt:      template.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      template.UpdatedAt.Format(time.RFC3339),
	}
}
