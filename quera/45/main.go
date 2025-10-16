package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func hello_get(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello, %s!", name)
}

func birthday(c *gin.Context) {
	name := c.Query("birthday")
	c.String(http.StatusOK, "Ha", name)
}
func main() {
	e := gin.Default()
	e.GET("/hello/:name", hello_get)

}
