package main

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
)

func TestEmptyData(t *testing.T) {
	data := map[string]string{}
	var body bytes.Buffer

	mw := multipart.NewWriter(&body)
	for key, value := range data {
		fw, err := mw.CreateFormField(key)
		assert.NoError(t, err)

		n, err := fw.Write([]byte(value))
		assert.NoError(t, err)
		assert.Equal(t, len(value), n)
	}
	err := mw.Close()
	assert.NoError(t, err)

	e := getEngine()
	req := httptest.NewRequest(http.MethodPost, "/user/validate", &body)
	req.Header.Set("Content-Type", binding.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())

	w := httptest.NewRecorder()

	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), `"message":"Data are invalid"`)
	h := w.Header()["Content-Type"]
	assert.Equal(t, len(h), 1)
	if len(h) == 0 {
		return
	}
	assert.Equal(t, h[0], "application/json; charset=utf-8")
}


func TestValidData(t *testing.T) {
	now := time.Now()
	birthDate := now.Add(-20 * 365 * 24 * time.Hour)
	data := map[string]string{
		"first_name": "ali",
		"last_name": "alavi",
		"username": "alialavi",
		"email": "alavi@quera.org",
		"phone_number": "09123456789",
		"birth_date": birthDate.Format("2006/01/02"),
		"national_id": "7421368515",
	}
	var body bytes.Buffer

	mw := multipart.NewWriter(&body)
	for key, value := range data {
		fw, err := mw.CreateFormField(key)
		assert.NoError(t, err)

		n, err := fw.Write([]byte(value))
		assert.NoError(t, err)
		assert.Equal(t, len(value), n)
	}
	err := mw.Close()
	assert.NoError(t, err)

	e := getEngine()
	req := httptest.NewRequest(http.MethodPost, "/user/validate", &body)
	req.Header.Set("Content-Type", binding.MIMEMultipartPOSTForm+"; boundary="+mw.Boundary())

	w := httptest.NewRecorder()
	
	e.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"message":"Data are valid"`)
	h := w.Header()["Content-Type"]
	assert.Equal(t, len(h), 1)
	if len(h) == 0 {
		return
	}
	assert.Equal(t, h[0], "application/json; charset=utf-8")
}