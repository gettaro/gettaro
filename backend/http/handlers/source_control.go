package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListSourceControlAccounts(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusOK, gin.H{"message": "List source control accounts endpoint"})
}

func CreateSourceControlAccount(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusCreated, gin.H{"message": "Create source control account endpoint"})
}
