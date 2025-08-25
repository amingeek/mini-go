package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	http.HandleFunc("/sayHelloWorld", handleHelloWorld)
	http.HandleFunc("/sayAny", sayAny)
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}

func root(w http.ResponseWriter, r *http.Request) {
	res := Response{
		Message: "Main Page",
		Status:  "OK",
	}
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func sayAny(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	res := &Response{}

	if text == "" {
		res = &Response{
			Message: "",
			Status:  "Bad Request",
		}
	} else {
		res = &Response{
			Message: text,
			Status:  "ok",
		}
	}
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Hello World!",
		Status:  "OK",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
