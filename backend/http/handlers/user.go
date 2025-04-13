package handlers

import (
	"net/http"

	"ems.dev/backend/http/middleware"
	userapi "ems.dev/backend/services/user/api"
	"ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userApi        *userapi.Api
	authMiddleware gin.HandlerFunc
}

func NewUserHandler(userApi *userapi.Api, authMiddleware gin.HandlerFunc) *UserHandler {
	return &UserHandler{
		userApi:        userApi,
		authMiddleware: middleware.AuthMiddleware(userApi),
	}
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/me", h.authMiddleware, h.GetMe)
}

// GetMe returns the current user's information
func (h *UserHandler) GetMe(c *gin.Context) {
	user, exists := c.Get("db_user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found in context"})
		return
	}

	dbUser, ok := user.(*types.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
		return
	}

	c.JSON(http.StatusOK, dbUser)
}
