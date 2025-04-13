package handlers

import (
	"net/http"

	"ems.dev/backend/http/middleware"
	httptypes "ems.dev/backend/http/types/user"
	usersApi "ems.dev/backend/services/user/api"
	"ems.dev/backend/services/user/types"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usersApi       *usersApi.Api
	authMiddleware gin.HandlerFunc
}

func NewUserHandler(usersApi *usersApi.Api) *UserHandler {
	return &UserHandler{
		usersApi:       usersApi,
		authMiddleware: middleware.AuthMiddleware(usersApi),
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/me", h.authMiddleware, h.GetMe)
}

// GetMe returns the current user's information
func (h *UserHandler) GetMe(c *gin.Context) {
	ctxUser, exists := c.Get("db_user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found in context"})
		return
	}

	user, ok := ctxUser.(*types.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
		return
	}

	foundUser, err := h.usersApi.FindUser(types.UserSearchParams{Email: &user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user: " + err.Error()})
		return
	}

	if foundUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, httptypes.GetUserResponse{User: foundUser})
}
