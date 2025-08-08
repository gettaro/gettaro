package utils

import (
	"errors"
	"net/http"

	liberrors "ems.dev/backend/libraries/errors"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var conflictErr *liberrors.ConflictError
	if errors.As(err, &conflictErr) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	var badRequestErr *liberrors.BadRequestError
	if errors.As(err, &badRequestErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
