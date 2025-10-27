package main

import "github.com/gin-gonic/gin"

func usernameLen(name string) bool {
	return len(name) >= 4
}

func checkPassword(password string, username string) bool {
	mirrorUsername := ""
	for i := len(username) - 1; i >= 0; i-- {
		mirrorUsername += string(username[i])
	}
	return password == mirrorUsername
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		password := c.GetHeader("password")

		if !usernameLen(username) || !checkPassword(password, username) {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
