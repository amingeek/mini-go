package main

import (
	"github.com/labstack/echo/v4"
)

func wsChatRoom(c echo.Context) error {
	// TODO
	return nil
}

func main() {
    e := echo.New()
	e.GET("/ws/chat/:roomId/user/:username", wsChatRoom)
	e.Logger.Fatal(e.Start(":8080"))
}