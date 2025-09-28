package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = "localhost:8080"

var upgrader = websocket.Upgrader{}

func wsHello(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	defer connection.Close()
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("recv: %s", message)
		response := []byte(fmt.Sprintf("Hello %s", message))
		err = connection.WriteMessage(messageType, response)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}
func main() {
	http.HandleFunc("/ws/hello", wsHello)
	log.Fatal(http.ListenAndServe(addr, nil))
}
