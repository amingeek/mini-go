package main

import (
	"fmt"
	"log"
	"time"

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

	go func() {
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("received message: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		tickerTime := <-ticker.C
		sendMessage := fmt.Sprintf("Quera college (time: %s)", tickerTime.String())
		err = connection.WriteMessage(websocket.TextMessage, []byte(sendMessage))
		if err != nil {
			log.Println("write error:", err)
		} else {
			log.Println("send message:", sendMessage)
		}
	}

}
