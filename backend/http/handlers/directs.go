package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateDirectReport(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create direct report endpoint"})
}
