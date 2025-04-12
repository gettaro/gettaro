package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateIntegrationConfig(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create integration config endpoint"})
}
