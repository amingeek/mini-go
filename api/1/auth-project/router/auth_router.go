package router

import (
	"auth-project/controller"
	"auth-project/helper"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controller.AuthController) *gin.Engine {
	router := gin.Default()

	// CORS middleware برای ارتباط با frontend
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(helper.AuthMiddleware())
	{
		protected.GET("/profile", func(ctx *gin.Context) {
			email := ctx.GetString("email")
			ctx.JSON(200, gin.H{
				"email":   email,
				"message": "This is a protected route",
			})
		})
	}

	return router
}
