package handlers

import (
	"net/http"

	httptypes "ems.dev/backend/http/types/title"
	"ems.dev/backend/http/utils"
	orgapi "ems.dev/backend/services/organization/api"
	titleapi "ems.dev/backend/services/title/api"
	titletypes "ems.dev/backend/services/title/types"
	"github.com/gin-gonic/gin"
)

// TitleHandler handles title-related HTTP requests
type TitleHandler struct {
	titleApi titleapi.TitleAPI
	orgApi   orgapi.OrganizationAPI
}

// NewTitleHandler creates a new TitleHandler instance
func NewTitleHandler(titleApi titleapi.TitleAPI, orgApi orgapi.OrganizationAPI) *TitleHandler {
	return &TitleHandler{
		titleApi: titleApi,
		orgApi:   orgApi,
	}
}

// CreateTitle handles the POST /api/organizations/{id}/titles endpoint
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 201: Title created successfully
// - 400: Bad request if organization ID is missing or request body is invalid
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Creates a new title in the database
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrDatabaseQuery: When database query fails
func (h *TitleHandler) CreateTitle(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req httptypes.CreateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title, err := h.titleApi.CreateTitle(c.Request.Context(), titletypes.Title{
		Name:           req.Name,
		OrganizationID: orgID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"title": title})
}

// ListTitles handles the GET /api/organizations/{id}/titles endpoint
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 200: List of titles
// - 400: Bad request if organization ID is missing
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Makes a database query to fetch all titles for the organization
// - Performs organization membership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrDatabaseQuery: When database query fails
func (h *TitleHandler) ListTitles(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is a member of the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	titles, err := h.titleApi.ListTitles(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"titles": titles})
}

// UpdateTitle handles the PUT /api/organizations/{id}/titles/{titleId} endpoint
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 200: Title updated successfully
// - 400: Bad request if organization ID or title ID is missing or request body is invalid
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 404: Not found if title does not exist
// - 500: Internal server error if service layer fails
// Side Effects:
// - Updates title details in the database
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrNotFound: When title is not found
// - ErrDatabaseQuery: When database query fails
func (h *TitleHandler) UpdateTitle(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	titleID := c.Param("titleId")
	if titleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req httptypes.UpdateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title, err := h.titleApi.UpdateTitle(c.Request.Context(), titleID, titletypes.Title{
		Name: req.Name,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if title == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "title not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"title": title})
}

// DeleteTitle handles the DELETE /api/organizations/{id}/titles/{titleId} endpoint
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 204: Title deleted successfully
// - 400: Bad request if organization ID or title ID is missing
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Deletes title from the database
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrDatabaseQuery: When database query fails
func (h *TitleHandler) DeleteTitle(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	titleID := c.Param("titleId")
	if titleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	if err := h.titleApi.DeleteTitle(c.Request.Context(), titleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignUserTitle handles the POST /api/organizations/{id}/titles/assign endpoint
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 201: User title assigned successfully
// - 400: Bad request if organization ID is missing or request body is invalid
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Assigns a title to a user in the database
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrDatabaseQuery: When database query fails
func (h *TitleHandler) AssignUserTitle(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req httptypes.AssignUserTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.titleApi.AssignUserTitle(c.Request.Context(), titletypes.UserTitle{
		UserID:         req.UserID,
		TitleID:        req.TitleID,
		OrganizationID: orgID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// RemoveUserTitle handles the DELETE /api/organizations/{id}/users/{userId}/title endpoint
// Params:
// - c: The Gin context containing request and response
// Returns:
// - 204: User title assignment removed successfully
// - 400: Bad request if organization ID or user ID is missing
// - 401: Unauthorized if user is not authenticated
// - 403: Forbidden if user does not have access to the organization
// - 500: Internal server error if service layer fails
// Side Effects:
// - Removes user title assignment from the database
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user does not have access to the organization
// - ErrDatabaseQuery: When database query fails
func (h *TitleHandler) RemoveUserTitle(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	if err := h.titleApi.RemoveUserTitle(c.Request.Context(), userID, orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registers all title-related routes
func (h *TitleHandler) RegisterRoutes(api *gin.RouterGroup) {
	organizations := api.Group("/organizations/:id")
	{
		// Title CRUD operations
		titles := organizations.Group("/titles")
		{
			titles.POST("", h.CreateTitle)
			titles.GET("", h.ListTitles)
			titles.PUT("/:titleId", h.UpdateTitle)
			titles.DELETE("/:titleId", h.DeleteTitle)
			titles.POST("/assign", h.AssignUserTitle)
		}

		// User title operations
		users := organizations.Group("/users/:userId")
		{
			users.DELETE("/title", h.RemoveUserTitle)
		}
	}
}
