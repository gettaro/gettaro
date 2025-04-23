package utils

import (
	"fmt"

	"ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(c *gin.Context) (*types.User, error) {
	ctxUser, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("user not found in context")
	}

	castedUser, ok := ctxUser.(*types.User)
	if !ok {
		return nil, fmt.Errorf("user is not of type types.User")
	}

	return castedUser, nil
}
