package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name        string `json:"name"`
	Family      string `json:"family"`
	Create_date time.Time
}

func WelcomeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome!")
}

func SignupHandler(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	user.Create_date = time.Now()
	db := initializeDB()
	db.Create(user)

	return c.JSON(http.StatusCreated, user)
}

func main() {
	e := echo.New()
	e.GET("/", WelcomeHandler)
	e.POST("/signup", SignupHandler)
	e.Start(":8080")
}
