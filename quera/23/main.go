package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	port int
	mux  *http.ServeMux
}

type APIResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func parseNumbers(param string) ([]int64, string) {
	if strings.TrimSpace(param) == "" {
		return nil, "'numbers' parameter missing"
	}
	parts := strings.Split(param, ",")
	nums := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, "invalid number"
		}
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, "invalid number"
		}
		nums = append(nums, n)
	}
	return nums, ""
}

func NewServer(port string) *Server {
	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	serverNew := Server{port: intPort}
	return &serverNew
}

func addChecked(a, b int64) (int64, bool) {
	if (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b) {
		return 0, false
	}
	return a + b, true

}

func subChecked(a, b int64) (int64, bool) {
	if b > 0 && a < math.MinInt64+b {
		return 0, false
	}
	if b < 0 && a > math.MaxInt64+b {
		return 0, false
	}
	return a - b, true
}

func add(w http.ResponseWriter, r *http.Request) {
	nums, errMsg := parseNumbers(r.URL.Query().Get("numbers"))
	if errMsg != "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{"", errMsg})
		return
	}

	var res int64 = 0
	for _, n := range nums {
		v, ok := addChecked(res, n)
		if !ok {
			writeJSON(w, http.StatusBadRequest, APIResponse{"", "Overflow"})
			return
		}
		res = v
	}

	writeJSON(w, http.StatusOK, APIResponse{fmt.Sprintf("The result of your query is: %d", res), ""})
}

func sub(w http.ResponseWriter, r *http.Request) {
	nums, errMsg := parseNumbers(r.URL.Query().Get("numbers"))
	if errMsg != "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{"", errMsg})
		return
	}

	if len(nums) == 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{"", "'numbers' parameter missing"})
		return
	}

	var res = nums[0]
	for _, n := range nums[1:] {
		v, ok := subChecked(res, n)
		if !ok {
			writeJSON(w, http.StatusBadRequest, APIResponse{"", "Overflow"})
			return
		}
		res = v
	}
	writeJSON(w, http.StatusOK, APIResponse{fmt.Sprintf("The result of your query is: %d", res), ""})
}

func (s *Server) Start() {
	http.HandleFunc("/add", add)
	http.HandleFunc("/sub", sub)
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
