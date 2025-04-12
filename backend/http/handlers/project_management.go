package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListPMAccounts(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusOK, gin.H{"message": "List project management accounts endpoint"})
}

func CreatePMAccount(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusCreated, gin.H{"message": "Create project management account endpoint"})
}

func ListPMTickets(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusOK, gin.H{"message": "List PM tickets endpoint"})
}

func CreatePMTicket(c *gin.Context) {
	// TODO: Implement authorization check
	c.JSON(http.StatusCreated, gin.H{"message": "Create PM ticket endpoint"})
}
