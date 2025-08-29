package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var requestData struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Borrow bool   `json:"borrow"`
}

type Book struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Borrowed bool   `json:"borrowed"`
}

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

func book(w http.ResponseWriter, r *http.Request) {
	Books := []Book{}

	switch r.Method {
	case "POST":
		title := strings.ToLower(r.URL.Query().Get("title"))
		author := strings.ToLower(r.URL.Query().Get("author"))
		if title == "" || author == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, i := range Books {
			if i.Title == title && i.Author == title {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		Books = append(Books, Book{
			Title:    title,
			Author:   author,
			Borrowed: false,
		})

	case "PUT":

	}

}

func (s *Server) Start() {

	http.HandleFunc("/book", book)
	port := fmt.Sprintf(":%d", s.port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
