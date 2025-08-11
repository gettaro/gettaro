package handlers

import (
	"net/http"

	"ems.dev/backend/http/types/member"
	"ems.dev/backend/http/utils"
	memberapi "ems.dev/backend/services/member/api"
	membertypes "ems.dev/backend/services/member/types"
	organizationapi "ems.dev/backend/services/organization/api"
	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	memberApi memberapi.MemberAPI
	orgApi    organizationapi.OrganizationAPI
}

func NewMemberHandler(memberApi memberapi.MemberAPI, orgApi organizationapi.OrganizationAPI) *MemberHandler {
	return &MemberHandler{
		memberApi: memberApi,
		orgApi:    orgApi,
	}
}

// ListOrganizationMembers handles listing all members of an organization
// It:
// 1. Validates the organization ID
// 2. Checks if the user has access to the organization
// 3. Retrieves all members of the organization
// Returns:
// - 200: List of organization members
// - 400: If the organization ID is missing
// - 401: If the user is not authenticated
// - 403: If the user does not have access to the organization
// - 500: If there's a database error
func (h *MemberHandler) ListOrganizationMembers(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user has access to the organization
	if !utils.CheckOrganizationMembership(c, h.orgApi, &orgID) {
		return
	}

	// Get organization members
	members, err := h.memberApi.GetOrganizationMembers(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// AddOrganizationMember handles adding a user as a member to an organization
// It:
// 1. Validates the organization ID
// 2. Checks if the current user is the owner
// 3. Validates the request body
// 4. Adds the specified user as a member
// Returns:
// - 201: If the member was added successfully
// - 400: If the request body is invalid
// - 401: If the user is not authenticated
// - 403: If the user is not the owner
// - 500: If there's a database error
func (h *MemberHandler) AddOrganizationMember(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req member.AddOrganizationMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.memberApi.AddOrganizationMember(c.Request.Context(), req.TitleID, req.SourceControlAccountID, &membertypes.OrganizationMember{
		OrganizationID: orgID,
		Email:          req.Email,
		Username:       req.Username,
	}); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// RemoveOrganizationMember handles removing a user from an organization
// It:
// 1. Validates the organization ID and member ID
// 2. Checks if the current user is the owner
// 3. Removes the specified member from the organization
// Returns:
// - 204: If the member was removed successfully
// - 400: If the IDs are missing
// - 401: If the user is not authenticated
// - 403: If the user is not the owner
// - 500: If there's a database error
func (h *MemberHandler) RemoveOrganizationMember(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberID := c.Param("memberId")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	// Get the member to find the user ID
	member, err := h.memberApi.GetOrganizationMemberByID(c.Request.Context(), memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if member == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	// Verify the member belongs to the specified organization
	if member.OrganizationID != orgID {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found in this organization"})
		return
	}

	if err := h.memberApi.RemoveOrganizationMember(c.Request.Context(), orgID, member.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateOrganizationMember handles updating a member's details in an organization
// It:
// 1. Validates the organization ID and member ID
// 2. Checks if the current user is the owner
// 3. Validates the request body
// 4. Updates the member's details (username, title, source control account)
// Returns:
// - 200: If the member was updated successfully
// - 400: If the request body is invalid or IDs are missing
// - 401: If the user is not authenticated
// - 403: If the user is not the owner
// - 404: If the member, title, or source control account is not found
// - 500: If there's a database error
// Side Effects:
// - Updates member details in the database
// - Performs organization ownership check
// Errors:
// - ErrMissingOrganizationID: When organization ID is missing from the request
// - ErrUnauthorized: When user is not authenticated
// - ErrForbidden: When user is not the owner
// - ErrNotFound: When member, title, or source control account is not found
// - ErrDatabaseQuery: When database query fails
func (h *MemberHandler) UpdateOrganizationMember(c *gin.Context) {
	orgID, err := utils.GetOrganizationIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	memberID := c.Param("memberId")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "member ID is required"})
		return
	}

	// Check if user is an owner of the organization
	if !utils.CheckOrganizationOwnership(c, h.orgApi, orgID) {
		return
	}

	var req member.UpdateOrganizationMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.memberApi.UpdateOrganizationMember(c.Request.Context(), orgID, memberID, req.TitleID, req.SourceControlAccountID, req.Username); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member updated successfully"})
}

// RegisterRoutes registers all member-related routes
func (h *MemberHandler) RegisterRoutes(api *gin.RouterGroup) {
	members := api.Group("/organizations/:id/members")
	{
		members.GET("", h.ListOrganizationMembers)
		members.POST("", h.AddOrganizationMember)
		members.PUT("/:memberId", h.UpdateOrganizationMember)
		members.DELETE("/:memberId", h.RemoveOrganizationMember)
	}
}
