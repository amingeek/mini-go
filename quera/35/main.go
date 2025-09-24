package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type JsonResult struct {
	Message string `json:"message"`
}

func SayHiHandler(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Name parameter is required")
	}
	res := &JsonResult{
		Message: "Hello " + name,
	}
	return c.JSON(http.StatusOK, res)
}

func main() {
	e := echo.New()
	e.GET("/sayhi", SayHiHandler)
	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
