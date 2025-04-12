package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListTeams(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List teams endpoint"})
}

func CreateTeam(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Create team endpoint"})
}

func AddTeamMember(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Add team member endpoint"})
}
