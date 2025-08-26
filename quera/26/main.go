package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port int
}

func NewServer(port string) *Server {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	return &Server{port: intPort}
}

func book(w http.ResponseWriter, r *http.Request) {}

func (s *Server) Start() {
	http.HandleFunc("/book", book)
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
