package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Borrowed bool   `json:"borrowed"`
}

type Server struct {
	port  int
	mu    sync.RWMutex
	books map[string]*Book
}

type resultErr struct {
	Result string `json:"Result"`
	Error  string `json:"Error"`
}

func NewServer(port string) *Server {
	p, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	return &Server{
		port:  p,
		books: make(map[string]*Book),
	}
}

func keyFor(title, author string) string {
	return strings.ToLower(strings.TrimSpace(title)) + "|" + strings.ToLower(strings.TrimSpace(author))
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) bookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}
		title := strings.ToLower(strings.TrimSpace(r.FormValue("title")))
		author := strings.ToLower(strings.TrimSpace(r.FormValue("author")))
		if title == "" || author == "" {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}

		k := keyFor(title, author)

		s.mu.Lock()
		if _, exists := s.books[k]; exists {
			s.mu.Unlock()
			writeJSON(w, http.StatusOK, resultErr{"this book is already in the library", ""})
			return
		}
		s.books[k] = &Book{Title: title, Author: author, Borrowed: false}
		s.mu.Unlock()

		writeJSON(w, http.StatusOK, resultErr{fmt.Sprintf("added book %s by %s", title, author), ""})
		return

	case http.MethodGet:
		title := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("title")))
		author := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("author")))
		if title == "" || author == "" {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}
		k := keyFor(title, author)

		s.mu.RLock()
		book, exists := s.books[k]
		s.mu.RUnlock()
		if !exists {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "this book does not exist"})
			return
		}
		if book.Borrowed {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "this book is borrowed"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"title": book.Title, "author": book.Author})
		return

	case http.MethodDelete:
		title := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("title")))
		author := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("author")))
		if title == "" || author == "" {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}
		k := keyFor(title, author)

		s.mu.Lock()
		_, exists := s.books[k]
		if !exists {
			s.mu.Unlock()
			writeJSON(w, http.StatusBadRequest, resultErr{"", "this book does not exist"})
			return
		}
		delete(s.books, k)
		s.mu.Unlock()
		writeJSON(w, http.StatusOK, resultErr{"successfully deleted", ""})
		return

	case http.MethodPut:
		title := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("title")))
		author := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("author")))
		if title == "" || author == "" {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}

		var body struct {
			Borrow *bool `json:"borrow"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Borrow == nil {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "borrow value cannot be empty"})
			return
		}

		k := keyFor(title, author)

		s.mu.Lock()
		book, exists := s.books[k]
		if !exists {
			s.mu.Unlock()
			writeJSON(w, http.StatusBadRequest, resultErr{"", "this book does not exist"})
			return
		}
		if *body.Borrow {
			if book.Borrowed {
				s.mu.Unlock()
				writeJSON(w, http.StatusBadRequest, resultErr{"", "this book is already borrowed"})
				return
			}
			book.Borrowed = true
			s.mu.Unlock()
			writeJSON(w, http.StatusOK, resultErr{"you have borrowed this book successfully", ""})
			return
		} else {
			if !book.Borrowed {
				s.mu.Unlock()
				writeJSON(w, http.StatusBadRequest, resultErr{"", "this book is already in the library"})
				return
			}
			book.Borrowed = false
			s.mu.Unlock()
			writeJSON(w, http.StatusOK, resultErr{"thank you for returning this book", ""})
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) Start() {
	http.HandleFunc("/book", s.bookHandler)
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
