package router

import (
	"api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	g1 := router.Group("/api")
	g1.Group("/api")
	{
		g1.GET("/ping/", handlers.Ping)
		g1.POST("/register/", handlers.Register)
		g1.GET("/login/", handlers.Login)
	}
}
