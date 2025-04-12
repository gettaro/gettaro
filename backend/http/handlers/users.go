package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List users endpoint"})
}

func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create user endpoint"})
}

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get user endpoint"})
}

func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update user endpoint"})
}

func DeleteUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
