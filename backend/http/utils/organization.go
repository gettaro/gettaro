package utils

import (
	"net/http"

	authTypes "ems.dev/backend/http/types/auth"
	orgapi "ems.dev/backend/services/organization/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

// CheckOrganizationOwnership checks if the authenticated user is an owner of the specified organization.
// If the user is not an owner, it returns false and sets the appropriate error response.
// If there's an error checking ownership, it returns false and sets the error response.
// If the user is an owner, it returns true.
func CheckOrganizationOwnership(c *gin.Context, orgApi orgapi.OrganizationAPI, orgID string) bool {
	// Get user claims from context
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return false
	}

	// Get user from database
	user, err := orgApi.(interface {
		FindUser(params usertypes.UserSearchParams) (*usertypes.User, error)
	}).FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return false
	}

	// Check if user is the owner of the organization
	isOwner, err := orgApi.IsOrganizationOwner(c.Request.Context(), orgID, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	if !isOwner {
		c.JSON(http.StatusForbidden, gin.H{"error": "only organization owners can perform this action"})
		return false
	}

	return true
}

// CheckOrganizationMembership checks if the authenticated user is a member of the specified organization.
// If the user is not a member, it returns false and sets the appropriate error response.
// If there's an error checking membership, it returns false and sets the error response.
// If the user is a member, it returns true.
func CheckOrganizationMembership(c *gin.Context, orgApi orgapi.OrganizationAPI, orgID *string) bool {
	if orgID == nil {
		return true // No organization ID specified, so no membership check needed
	}

	// Get user claims from context
	userClaims, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return false
	}

	// Get user from database
	user, err := orgApi.(interface {
		FindUser(params usertypes.UserSearchParams) (*usertypes.User, error)
	}).FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return false
	}

	// Get user's organizations
	orgs, err := orgApi.GetUserOrganizations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	// Check if user is a member of the specified organization
	for _, org := range orgs {
		if org.ID == *orgID {
			return true
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this organization"})
	return false
}
