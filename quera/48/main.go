package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func getEngine() *gin.Engine {
	e := gin.Default()
	e.Use(authMiddleware())

	e.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	return e
}

func main() {
	e := getEngine()
	e.Run()
}