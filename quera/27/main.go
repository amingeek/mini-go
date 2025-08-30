package main

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port string
}

func NewServer(port string) *Server {
	return &Server{Port: port}
}

func (server *Server) startServer() {
	err := http.ListenAndServe(server.Port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server started at port %d\n", server.Port)
}

func main() {
	server := NewServer(":8000")
	server.startServer()
}
