package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func validateUser(c *gin.Context) {
	user := new (User)
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Data are invalid",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data are valid",
		"user": user,
	})
}

func getEngine() *gin.Engine {
	e := gin.Default()
	registerValidators()
	e.POST("/user/validate", validateUser)
	return e
}

func main() {
	e := getEngine()
	e.Run()
}