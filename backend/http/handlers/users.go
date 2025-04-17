package handlers

import (
	"net/http"

	"ems.dev/backend/http/middleware"
	userhttptypes "ems.dev/backend/http/types/user"
	usersApi "ems.dev/backend/services/user/api"
	"ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	usersApi       *usersApi.Api
	authMiddleware gin.HandlerFunc
}

func NewUsersHandler(usersApi *usersApi.Api) *UsersHandler {
	return &UsersHandler{
		usersApi:       usersApi,
		authMiddleware: middleware.AuthMiddleware(usersApi),
	}
}

func (h *UsersHandler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users endpoint"})
}

func (h *UsersHandler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create user endpoint"})
}

func (h *UsersHandler) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user endpoint"})
}

func (h *UsersHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user endpoint"})
}

func (h *UsersHandler) DeleteUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// GetMe returns the current user's information
func (h *UsersHandler) GetMe(c *gin.Context) {
	ctxUser, exists := c.Get("user")
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

	c.JSON(http.StatusOK, userhttptypes.GetUserResponse{User: foundUser})
}
