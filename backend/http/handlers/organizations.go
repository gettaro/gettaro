package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrganization(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create organization endpoint"})
}

func ListOrganizations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List organizations endpoint"})
}
