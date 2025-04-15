package handlers

import (
	"net/http"
	"strings"

	authTypes "ems.dev/backend/http/types/auth"
	orgapi "ems.dev/backend/services/organization/api"
	orgtypes "ems.dev/backend/services/organization/types"
	userapi "ems.dev/backend/services/user/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgApi  *orgapi.Api
	userApi *userapi.Api
}

func NewOrganizationHandler(orgApi *orgapi.Api, userApi *userapi.Api) *OrganizationHandler {
	return &OrganizationHandler{
		orgApi:  orgApi,
		userApi: userApi,
	}
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req orgtypes.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user_claims from context (set by auth middleware)
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, err := h.userApi.FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"organization": org})
}

func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	// Get user_claims from context (set by auth middleware)
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, err := h.userApi.FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
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

func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	// Get user_claims from context (set by auth middleware)
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, err := h.userApi.FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
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

func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	// Get user_claims from context (set by auth middleware)
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, err := h.userApi.FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
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

func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization ID is required"})
		return
	}

	err := h.orgApi.DeleteOrganization(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
