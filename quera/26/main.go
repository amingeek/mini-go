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
	port  int
	mutex sync.RWMutex
	books map[string]*Book
}
type resultErr struct {
	Result string `json:"Result"`
	Error  string `json:"Error"`
}

func NewServer(port string) *Server {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	return &Server{port: intPort, books: make(map[string]*Book)}

}
func keyFor(title, author string) string {
	return strings.ToLower(strings.TrimSpace(title)) + "|" + strings.ToLower(strings.TrimSpace(author))
}
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) book(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		// داده‌ها به صورت فرم ارسال می‌شوند
		if err := r.ParseForm(); err != nil {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}
		title := strings.TrimSpace(r.FormValue("title"))
		author := strings.TrimSpace(r.FormValue("author"))
		if title == "" || author == "" {
			writeJSON(w, http.StatusBadRequest, resultErr{"", "title or author cannot be empty"})
			return
		}

		k := keyFor(title, author)

		s.mutex.Lock()
		defer s.mutex.Unlock()
		if _, exists := s.books[k]; exists {
			// اگر قبلاً اضافه شده باشد
			writeJSON(w, http.StatusOK, resultErr{"this book is already in the library", ""})
			return
		}

		// اضافه کردن کتاب با حفظ حالت حروف اصلی (همان‌طور که فرستاده شده)
		s.books[k] = &Book{
			Title:    title,
			Author:   author,
			Borrowed: false,
		}
		writeJSON(w, http.StatusOK, resultErr{fmt.Sprintf("added book %s by %s", title, author), ""})
		return
	case "PUT":

	}

}

func (s *Server) Start() {

	http.HandleFunc("/book", s.book)
	port := fmt.Sprintf(":%d", s.port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
