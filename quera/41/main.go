package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	url := "ws://localhost:8080/ws/hello"
	log.Printf("connecting to %s", url)

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer connection.Close()

	sendMessage := "Quera college"
	err = connection.WriteMessage(websocket.TextMessage, []byte(sendMessage))
	if err != nil {
		log.Println("write error:", err)
	} else {
		log.Println("send message:", sendMessage)
	}
	_, message, err := connection.ReadMessage()
	if err != nil {
		log.Println("read error:", err)
		return
	}
	log.Printf("received message: %s", message)
}
