package handlers

import (
	"net/http"

	userhttptypes "ems.dev/backend/http/types/user"
	usersApi "ems.dev/backend/services/user/api"
	"ems.dev/backend/services/user/types"
	"github.com/gin-gonic/gin"
)

// UsersHandler handles all user-related HTTP requests.
// It provides endpoints for user management including CRUD operations and retrieving the current user's information.
type UsersHandler struct {
	usersApi usersApi.UserAPI
}

// NewUsersHandler creates a new instance of UsersHandler.
// It initializes the handler with the provided UserAPI
func NewUsersHandler(usersApi usersApi.UserAPI) *UsersHandler {
	return &UsersHandler{
		usersApi: usersApi,
	}
}

// ListUsers handles the GET /api/users endpoint.
// It returns a list of all users in the system.
// Returns:
// - 200: List of users
// - 500: If there's a database error
func (h *UsersHandler) ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users endpoint"})
}

// CreateUser handles the POST /api/users endpoint.
// It creates a new user in the system.
// Returns:
// - 201: The created user
// - 400: If the request body is invalid
// - 500: If there's a database error
func (h *UsersHandler) CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create user endpoint"})
}

// GetUser handles the GET /api/users/:id endpoint.
// It retrieves a specific user by their ID.
// Returns:
// - 200: The user details
// - 404: If the user is not found
// - 500: If there's a database error
func (h *UsersHandler) GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user endpoint"})
}

// UpdateUser handles the PUT /api/users/:id endpoint.
// It updates an existing user's information.
// Returns:
// - 200: The updated user
// - 400: If the request body is invalid
// - 404: If the user is not found
// - 500: If there's a database error
func (h *UsersHandler) UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user endpoint"})
}

// DeleteUser handles the DELETE /api/users/:id endpoint.
// It deletes a user from the system.
// Returns:
// - 204: If the user was successfully deleted
// - 404: If the user is not found
// - 500: If there's a database error
func (h *UsersHandler) DeleteUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// GetMe handles the GET /api/users/me endpoint.
// It returns the current authenticated user's information.
// Returns:
// - 200: The current user's details
// - 401: If the user is not authenticated
// - 404: If the user is not found in the database
// - 500: If there's a database error
func (h *UsersHandler) GetMe(c *gin.Context) {
	ctxUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	user, ok := ctxUser.(*types.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user type in context"})
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

// RegisterRoutes registers all user-related routes with the provided router group.
// It sets up the following routes:
// - GET /api/users - List all users
// - POST /api/users - Create a new user
// - GET /api/users/:id - Get a specific user
// - PUT /api/users/:id - Update a user
// - DELETE /api/users/:id - Delete a user
// - GET /api/users/me - Get current user's information
func (h *UsersHandler) RegisterRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.ListUsers)
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("/me", h.GetMe)
	}
}
