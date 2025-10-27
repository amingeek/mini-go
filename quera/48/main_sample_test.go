package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyHeaders(t *testing.T) {
	e := getEngine()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Contains(t, w.Body.String(), `"message":"Unauthorized"`)
	h := w.Header()["Content-Type"]
	assert.Equal(t, len(h), 1)
	if len(h) == 0 {
		return
	}
	assert.Equal(t, h[0], "application/json; charset=utf-8")
}

func TestInvalidUserPass(t *testing.T) {
	e := getEngine()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("username", "abcd1234")
	req.Header.Set("password", "abcd1234")

	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusUnauthorized)
	assert.Contains(t, w.Body.String(), `"message":"Unauthorized"`)
	h := w.Header()["Content-Type"]
	assert.Equal(t, len(h), 1)
	if len(h) == 0 {
		return
	}
	assert.Equal(t, h[0], "application/json; charset=utf-8")
}

func TestValidUserPass(t *testing.T) {
	e := getEngine()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("username", "abcd1234")
	req.Header.Set("password", "4321dcba")

	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Contains(t, w.Body.String(), `pong`)
	h := w.Header()["Content-Type"]
	assert.Equal(t, len(h), 1)
	if len(h) == 0 {
		return
	}
	assert.Equal(t, h[0], "text/plain; charset=utf-8")
}

