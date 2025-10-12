package main

import (
	_ "encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Message struct {
	User  string `json:"user"`
	Text  string `json:"text"`
	Color string `json:"color"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/ws", func(c echo.Context) error {
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()
		clients[ws] = true

		for {
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				delete(clients, ws)
				break
			}
			broadcast <- msg
		}
		return nil
	})

	go handleMessages()

	fmt.Println("✅ WebSocket server started on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
