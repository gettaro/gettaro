package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListPullRequests(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusOK, gin.H{"message": "List pull requests endpoint"})
}

func CreatePullRequest(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusCreated, gin.H{"message": "Create pull request endpoint"})
}

func CreatePRComment(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusCreated, gin.H{"message": "Create PR comment endpoint"})
}

func AddPRReviewer(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusCreated, gin.H{"message": "Add PR reviewer endpoint"})
}
