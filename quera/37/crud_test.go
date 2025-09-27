package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	productsId = []uint{}
)

func parseResponse(t *testing.T, rec *httptest.ResponseRecorder) ApiResponse {
	var resp ApiResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	return resp
}

func extractProduct(t *testing.T, data interface{}) Product {
	dataBytes, _ := json.Marshal(data)
	var p Product
	json.Unmarshal(dataBytes, &p)
	return p
}

func TestCreateProduct(t *testing.T) {
	e := echo.New()
	db = initDB()

	product := Product{
		Name:        "Test Laptop Dell",
		Price:       2000,
		Color:       "Silver",
		Category:    "Laptop",
		CreatedDate: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	body, _ := json.Marshal(product)

	req := httptest.NewRequest(http.MethodPost, "/products/", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, createProduct(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Empty(t, resp.Error)

		p := extractProduct(t, resp.Data)
		assert.Equal(t, "Test Laptop Dell", p.Name)
		productsId = append(productsId, p.ID)
	}
}

func TestGetProducts(t *testing.T) {
	e := echo.New()
	db = initDB()

	req := httptest.NewRequest(http.MethodGet, "/products/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, getProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Empty(t, resp.Error)

		dataBytes, _ := json.Marshal(resp.Data)
		var products []Product
		json.Unmarshal(dataBytes, &products)
		assert.GreaterOrEqual(t, len(products), 1)
	}
}

func TestGetCategory(t *testing.T) {
	e := echo.New()
	db = initDB()

	category := "Laptop"
	req := httptest.NewRequest(http.MethodGet, "/products/category/"+category, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("category")
	c.SetParamValues(category)

	if assert.NoError(t, getCategory(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Empty(t, resp.Error)

		dataBytes, _ := json.Marshal(resp.Data)
		var products []Product
		json.Unmarshal(dataBytes, &products)
		for _, p := range products {
			assert.True(t, strings.EqualFold(category, p.Category))
		}
	}
}

func TestGetProduct(t *testing.T) {
	e := echo.New()
	db = initDB()

	product := Product{
		Name:        "Test Phone",
		Price:       999,
		Color:       "Blue",
		Category:    "Mobile",
		CreatedDate: time.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC),
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

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Empty(t, resp.Error)

		p := extractProduct(t, resp.Data)
		assert.Equal(t, "Test Phone", p.Name)
		productsId = append(productsId, p.ID)
	}
}

func TestDeleteProduct(t *testing.T) {
	e := echo.New()
	db = initDB()

	product := Product{
		Name:        "Test Watch Delete",
		Price:       2000,
		Color:       "Black",
		Category:    "Watch",
		CreatedDate: time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("failed to insert product: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/"+fmt.Sprint(product.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(product.ID))

	if assert.NoError(t, deleteProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Empty(t, resp.Error)

		p := extractProduct(t, resp.Data)
		assert.Equal(t, "Test Watch Delete", p.Name)
	}
}

func TestDeleteProducts(t *testing.T) {
	if assert.NoError(t, deleteProducts(productsId)) {
		assert.Equal(t, nil, nil)
	}
}
