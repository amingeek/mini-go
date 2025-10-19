package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	e := getEngine()
	req := httptest.NewRequest(http.MethodPost, "/register", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, w.Body.String(), `{"message":"firtname is required"}`)
	h := w.Header()["Content-Type"]
	assert.Equal(t, len(h), 1)
	if len(h) == 0 {
		return
	}
	assert.Equal(t, h[0], "application/json; charset=utf-8")
}


func TestHello(t *testing.T) {
	e := getEngine()
	req := httptest.NewRequest(http.MethodGet, "/hello/ali/alavi", nil)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Body.String(), `ali alavi is not registered`)
	h := w.Header()["Content-Type"]
	if len(h) == 0 {
		return
	}
	assert.Equal(t, len(h), 1)
	assert.Equal(t, h[0], "text/plain; charset=utf-8")
}