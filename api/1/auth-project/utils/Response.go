package utils

import (
	"api/models"
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, models.Response{
		Success: true,
		Message: message,
		Data:    data,
		Code:    statusCode,
	})
}
