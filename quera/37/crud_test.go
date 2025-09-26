package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateProductWithRealDB(t *testing.T) {
	e := echo.New()
	db = initDB()

	product := Product{
		Name:        "Laptop Dell",
		Price:       2000,
		Color:       "Silver",
		Category:    "Laptop",
		CreatedDate: time.Now(),
	}
	body, _ := json.Marshal(product)

	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, createProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "Laptop Dell")
		fmt.Println("TestCreateProductWithRealDB PASS")
	}
}

func TestGetProducts(t *testing.T) {
	e := echo.New()
	db = initDB()

	product := Product{
		Name:        "Test Laptop",
		Price:       1500,
		Color:       "Black",
		Category:    "Laptop",
		CreatedDate: time.Now(),
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to insert product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/products/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, getProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Test Laptop")
		fmt.Println("TestGetProducts PASS")
	}
}

func TestGetCategory(t *testing.T) {
	e := echo.New()
	db = initDB()
	product := Product{
		Name:        "Test Laptop Category",
		Price:       1500,
		Color:       "Black",
		Category:    "Laptop",
		CreatedDate: time.Now(),
	}
	// write tomorrow
}

func TestGetProductWithRealDB(t *testing.T) {
	e := echo.New()

	db = initDB()

	product := Product{
		Name:        "Test Phone",
		Price:       999,
		Color:       "Blue",
		Category:    "Mobile",
		CreatedDate: time.Now(),
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to insert product: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/products/"+fmt.Sprint(product.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(product.ID))

	if assert.NoError(t, getProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Test Phone")
		fmt.Println("TestGetProductWithRealDB PASS")
	}
}
