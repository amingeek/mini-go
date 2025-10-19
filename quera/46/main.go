package main

import "github.com/gin-gonic/gin"

func getEngine() *gin.Engine {
	e := gin.Default()
	e.POST("/register", register)
	e.GET("/hello/:firstname/:lastname", hello)
	return e
}

func main() {
	e := getEngine()
	e.Run()
}
