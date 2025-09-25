package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func testCreateProduct(t *testing.T) {

	e := echo.New()

	t.Run("craeteProduct", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// call handler
		if assert.NoError(t, SayHiHandler(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "Name parameter is required", rec.Body.String())
		}
	})
}
