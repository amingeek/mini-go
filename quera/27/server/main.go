package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
)

type Server struct {
	Port string
}

type resultErr struct {
	Result string `json:"Result"`
	Error  string `json:"Error"`
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func NewServer(port string) *Server {
	return &Server{Port: port}
}

func getWeatherData() string {
	err := godotenv.Load()
	if err != nil {
		return fmt.Sprintf("Error loading .env file")
	}
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		return fmt.Sprintf("OPENWEATHERMAP_API_KEY is not set in the .env file")
	}
	apiUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=Tehran&appid=%s", apiKey)

	response, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Sprintf("Error making the API request:", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Sprintf("Error reading the response body:", err)
	}

	return fmt.Sprintf(string(responseBody))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := getWeatherData()
	res := resultErr{
		Result: data,
		Error:  "",
	}
	writeJson(w, http.StatusOK, res)
	return

}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		writeJson(w, http.StatusOK, resultErr{"", "name is empty"})
		return
	}
	res := resultErr{
		Result: fmt.Sprintf("Hello, %s!", name),
		Error:  "",
	}
	writeJson(w, http.StatusOK, res)
	return
}
func (server *Server) Start() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/json", jsonHandler)
	err := http.ListenAndServe(server.Port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server started at port %d\n", server.Port)
}

func main() {
	server := NewServer(":8000")
	server.Start()
}
