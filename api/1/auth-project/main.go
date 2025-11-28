package main

import (
	"api/database"
	"api/routes"
	"api/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ServerPort = ":8080"
)

func main() {
	database.Connect()

	fmt.Println("Database connected 🆗")

	r := gin.Default()
	router.SetupRoutes(r)

	fmt.Println(fmt.Sprintf("Server Started port %s🆗", ServerPort))
	psd, _ := utils.PasswordHash("1388")
	fmt.Println(psd)
	r.Run(ServerPort)
}
