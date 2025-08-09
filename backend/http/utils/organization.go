package utils

import (
	"fmt"
	"net/http"

	orgapi "ems.dev/backend/services/organization/api"
	usertypes "ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

// GetOrganizationIDFromContext extracts the organization ID from the request context and returns it
// It validates that the ID parameter is present in the URL
// Returns:
// - string: The organization ID if present
// - error: If the ID parameter is missing
func GetOrganizationIDFromContext(c *gin.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("organization ID is required")
	}
	return id, nil
}

// CheckOrganizationOwnership checks if the authenticated user is an owner of the specified organization.
// If the user is not an owner, it returns false and sets the appropriate error response.
// If there's an error checking ownership, it returns false and sets the error response.
// If the user is an owner, it returns true.
func CheckOrganizationOwnership(c *gin.Context, orgApi orgapi.OrganizationAPI, orgID string) bool {
	// Get user claims from context
	ctxUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return false
	}

	user := ctxUser.(*usertypes.User)

	// Get get user organizations
	userOrganizations, err := orgApi.GetMemberOrganizations(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	// Check if user is the owner of the organization
	for _, org := range userOrganizations {
		if org.ID == orgID && org.IsOwner {
			return true
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "only organization owners can perform this action"})
	return false
}

// CheckOrganizationMembership checks if the authenticated user is a member of the specified organization.
// If the user is not a member, it returns false and sets the appropriate error response.
// If there's an error checking membership, it returns false and sets the error response.
// If the user is a member, it returns true.
func CheckOrganizationMembership(c *gin.Context, orgApi orgapi.OrganizationAPI, orgID *string) bool {
	if orgID == nil {
		return true // No organization ID specified, so no membership check needed
	}

	// Get user from context
	ctxUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return false
	}

	// Get user from database
	userOrganizations, err := orgApi.GetMemberOrganizations(c.Request.Context(), ctxUser.(*usertypes.User).ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}

	for _, org := range userOrganizations {
		if org.ID == *orgID {
			return true
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "user is not a member of this organization"})
	return false
}
