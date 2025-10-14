package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "hello %s", name)
}

func main() {
	e := gin.Default()
	e.GET("/hello/:name", hello)
	e.Run()
}
