package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSayHiHandler(t *testing.T) {
	e := echo.New()

	t.Run("missing name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sayhi", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// call handler
		if assert.NoError(t, SayHiHandler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Name parameter is required", rec.Body.String())
		}
	})

	t.Run("with name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sayhi?name=Amir", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// call handler
		if assert.NoError(t, SayHiHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"message":"Hello Amir"}`, rec.Body.String())
		}
	})

	t.Run("with other name", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sayhi?name=Amin", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, SayHiHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"message":"Hello Amin"}`, rec.Body.String())
		}
	})

	t.Run("other address", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		if assert.NoError(t, SayHiHandler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}
