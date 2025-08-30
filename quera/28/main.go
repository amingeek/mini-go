package main

import (
	// "net/http"
	"fmt"
	"net/http"
)

type Server struct {
	Port string `json:"Port"`
}

func NewServer(port string) *Server {
	return &Server{Port: port}
}

func (s *Server) Start() {
	err := http.ListenAndServe(":"+s.Port, nil)
	if err != nil {
		fmt.Println("Error sending request:", err)
	} else {
		fmt.Println("Successfully started server")
	}
}

func main() {
	server := NewServer("8080")
	server.Start()
}
