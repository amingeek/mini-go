package main

import (
	"TestEcho/database"
	"TestEcho/handlers"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func WithDBConnection(handlerFunc func(c echo.Context, db *gorm.DB) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		dbConn, err := database.GetConnection()
		if err != nil {
			return c.JSON(http.StatusBadGateway, "Can't Connect To Database")
		}
		return handlerFunc(c, dbConn)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	err = database.Connect()
	if err != nil {
		fmt.Println("Failed to connecting to database")
		return
	}
	e := echo.New()
	e.GET("/products", WithDBConnection(handlers.ListProducts))
	e.POST("/add-product", WithDBConnection(handlers.AddProduct))
	log.Fatal(e.Start(":8080"))
}
