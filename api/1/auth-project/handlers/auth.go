package handlers

import (
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func register(c *gin.Context) {
	var req models.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}


}
