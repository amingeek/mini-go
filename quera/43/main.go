package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rooms = make(map[string][]*websocket.Conn)
	mutex sync.Mutex
)

func wsChatRoom(c echo.Context) error {
	roomID := c.Param("roomId")
	username := c.Param("username") // اینجا باید استفاده بشه

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	mutex.Lock()
	rooms[roomID] = append(rooms[roomID], conn)
	mutex.Unlock()

	go func() {
		defer func() {
			conn.Close()
			mutex.Lock()
			conns := rooms[roomID]
			for i, c := range conns {
				if c == conn {
					rooms[roomID] = append(conns[:i], conns[i+1:]...)
					break
				}
			}
			if len(rooms[roomID]) == 0 {
				delete(rooms, roomID)
			}
			mutex.Unlock()
		}()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			finalMsg := username + ": " + string(msg)

			mutex.Lock()
			for _, otherConn := range rooms[roomID] {
				if otherConn != conn {
					otherConn.WriteMessage(websocket.TextMessage, []byte(finalMsg))
				}
			}
			mutex.Unlock()
		}
	}()

	return nil
}

func main() {
	e := echo.New()
	e.GET("/ws/chat/:roomId/user/:username", wsChatRoom)
	e.Logger.Fatal(e.Start(":8080"))
}
