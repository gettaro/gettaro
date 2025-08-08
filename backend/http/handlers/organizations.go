package handlers

import (
	"fmt"
	"net/http"
	"strings"

	orghttptypes "ems.dev/backend/http/types/organization"
	"ems.dev/backend/http/utils"
	"ems.dev/backend/services/errors"
	orgapi "ems.dev/backend/services/organization/api"
	orgtypes "ems.dev/backend/services/organization/types"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgApi  orgapi.OrganizationAPI
	userApi userapi.UserAPI
}

func NewOrganizationHandler(orgApi orgapi.OrganizationAPI, userApi userapi.UserAPI) *OrganizationHandler {
	return &OrganizationHandler{
		orgApi:  orgApi,
		userApi: userApi,
	}
}

// getUserFromContext extracts the user from the request context and returns it
// It retrieves the user claims from the JWT token and looks up the corresponding user in the database
// Returns:
// - *usertypes.User: The user object if found
// - error: If the user claims are missing, user not found, or database error occurs
func (h *OrganizationHandler) getUserFromContext(c *gin.Context) (*usertypes.User, error) {
	ctxUser, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	castedUser, ok := ctxUser.(*usertypes.User)
	if !ok {
		return nil, fmt.Errorf("user is not of type usertypes.User")
	}

	return castedUser, nil
}

// getOrganizationIDFromContext extracts the organization ID from the request context and returns it
// It validates that the ID parameter is present in the URL
// Returns:
// - string: The organization ID if present
// - error: If the ID parameter is missing
func (h *OrganizationHandler) getOrganizationIDFromContext(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("organization ID is required")
	}
	return id, nil
}

// CreateOrganization handles the creation of a new organization
// It:
// 1. Validates the request body
// 2. Gets the current user from context
// 3. Creates a new organization with the provided name and slug
// 4. Sets the current user as the owner
// Returns:
// - 201: The created organization
// - 400: If the request body is invalid
// - 401: If the user is not authenticated
// - 409: If the organization slug already exists
// - 500: If there's a database error
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req orgtypes.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Create organization
	org := orgtypes.Organization{
		Name: req.Name,
		Slug: strings.ToLower(req.Slug),
	}

	// Create organization and set user as owner
	err = h.orgApi.CreateOrganization(c.Request.Context(), &org, user.ID)
	if err != nil {
		if errors.IsDuplicateConflict(err) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"organization": org})
}

// ListOrganizations handles listing all organizations the current user is a member of
// It:
// 1. Gets the current user from context
// 2. Retrieves all organizations where the user is a member
// Returns:
// - 200: List of organizations with ownership information
// - 401: If the user is not authenticated
// - 500: If there's a database error
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get user's organizations
	orgs, err := h.orgApi.GetUserOrganizations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organizations": orgs})
}

// GetOrganization handles retrieving a specific organization by ID
// It:
// 1. Gets the current user from context
// 2. Validates the organization ID
// 3. Checks if the user has access to the organization
// Returns:
// - 200: The organization details
// - 400: If the organization ID is missing
// - 401: If the user is not authenticated
// - 404: If the organization is not found or user has no access
// - 500: If there's a database error
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get organization and check if user has access
	orgs, err := h.orgApi.GetUserOrganizations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var org *orgtypes.Organization
	for _, o := range orgs {
		if o.ID == id {
			org = &o
			break
		}
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organization": org})
}

// UpdateOrganization handles updating an existing organization
// It:
// 1. Gets the current user from context
// 2. Validates the organization ID
// 3. Checks if the user is the owner of the organization
// 4. Updates the organization with the provided fields
// Returns:
// - 200: The updated organization
// - 400: If the request body is invalid
// - 401: If the user is not authenticated
// - 403: If the user is not the owner
// - 404: If the organization is not found
// - 500: If there's a database error
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	user, err := h.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req orgtypes.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get organization and check if user is owner
	orgs, err := h.orgApi.GetUserOrganizations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var org *orgtypes.Organization
	for _, o := range orgs {
		if o.ID == id && o.IsOwner {
			org = &o
			break
		}
	}

	if org == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found or user is not the owner"})
		return
	}

	// Update fields
	if req.Name != "" {
		org.Name = req.Name
	}
	if req.Slug != "" {
		org.Slug = strings.ToLower(req.Slug)
	}

	// Save changes
	err = h.orgApi.UpdateOrganization(c.Request.Context(), org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"organization": org})
}

// DeleteOrganization handles deleting an organization
// It:
// 1. Validates the organization ID
// 2. Deletes the organization
// Returns:
// - 204: If the organization was deleted successfully
// - 400: If the organization ID is missing
// - 500: If there's a database error
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, id) {
		return
	}

	err := h.orgApi.DeleteOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddOrganizationMember handles adding a user as a member to an organization
// It:
// 1. Gets the current user from context
// 2. Validates the organization ID
// 3. Checks if the current user is the owner
// 4. Adds the specified user as a member
// Returns:
// - 201: If the member was added successfully
// - 400: If the request body is invalid
// - 401: If the user is not authenticated
// - 403: If the user is not the owner
// - 500: If there's a database error
func (h *OrganizationHandler) AddOrganizationMember(c *gin.Context) {
	orgID := c.Param("id")
	if orgID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req orghttptypes.AddOrganizationMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orgApi.AddOrganizationMember(c.Request.Context(), &orgtypes.UserOrganization{
		OrganizationID: orgID,
		Email:          req.Email,
		Username:       req.Username,
	}); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

// RemoveOrganizationMember handles removing a user from an organization
// It:
// 1. Gets the current user from context
// 2. Validates the organization ID and user ID
// 3. Checks if the current user is the owner
// 4. Removes the specified user from the organization
// Returns:
// - 204: If the member was removed successfully
// - 400: If the IDs are missing
// - 401: If the user is not authenticated
// - 403: If the user is not the owner
// - 500: If there's a database error
func (h *OrganizationHandler) RemoveOrganizationMember(c *gin.Context) {
	orgID := c.Param("id")
	userID := c.Param("userId")
	if orgID == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID and user ID are required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	if err := h.orgApi.RemoveOrganizationMember(c.Request.Context(), orgID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListOrganizationMembers handles listing all members of an organization
// It:
// 1. Validates the organization ID
// 2. Retrieves all members of the organization
// Returns:
// - 200: List of organization members
// - 400: If the organization ID is missing
// - 500: If there's a database error
func (h *OrganizationHandler) ListOrganizationMembers(c *gin.Context) {
	orgID, err := h.getOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	// Get organization members
	members, err := h.orgApi.GetOrganizationMembers(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// RegisterRoutes registers all organization-related routes
func (h *OrganizationHandler) RegisterRoutes(api *gin.RouterGroup) {
	organizations := api.Group("/organizations")
	{
		// Organization CRUD operations
		organizations.POST("", h.CreateOrganization)
		organizations.GET("", h.ListOrganizations)
		organizations.GET("/:id", h.GetOrganization)
		organizations.PUT("/:id", h.UpdateOrganization)
		organizations.DELETE("/:id", h.DeleteOrganization)

		// Organization member operations
		organizations.GET("/:id/members", h.ListOrganizationMembers)
		organizations.POST("/:id/members", h.AddOrganizationMember)
		organizations.DELETE("/:id/members/:userId", h.RemoveOrganizationMember)
	}
}
