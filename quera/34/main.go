package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func WelcomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome!")
}

func main() {
	e := echo.New()
	e.GET("/", WelcomeHandler)
	e.Start(":8080")
}
