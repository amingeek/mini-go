package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Count       int       `json:"count"`
	Price       float64   `json:"price"`
	Color       string    `json:"color"`
	Category    string    `json:"category"`
	CreatedDate time.Time `json:"created_date"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var (
	db = initDB()
)

func createProduct(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Error:   "invalid request body",
		})
	}

	validations := map[string]bool{
		"name is required":         product.Name == "",
		"count is required":        product.Count == 0,
		"price is required":        product.Price == 0,
		"color is required":        product.Color == "",
		"category is required":     product.Category == "",
		"created_date is required": product.CreatedDate.IsZero(),
	}
	for msg, failed := range validations {
		if failed {
			return c.JSON(http.StatusBadRequest, ApiResponse{
				Success: false,
				Error:   msg,
			})
		}
	}

	var existing Product
	if err := db.Where("name = ?", product.Name).First(&existing).Error; err == nil {
		return c.JSON(http.StatusBadRequest, ApiResponse{
			Success: false,
			Error:   "product with this name already exists",
		})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	if err := db.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, ApiResponse{
		Success: true,
		Data:    product,
	})
}

func getProducts(c echo.Context) error {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, ApiResponse{Success: true, Data: products})
}

func getCategory(c echo.Context) error {
	var products []Product
	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, ApiResponse{Success: false, Error: "category is required"})
	}
	if err := db.Where("category = ?", category).Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, ApiResponse{Success: true, Data: products})
}

func getProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, ApiResponse{Success: false, Error: "id is required"})
	}

	var product Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, ApiResponse{Success: false, Error: "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}
	return c.JSON(http.StatusOK, ApiResponse{Success: true, Data: product})
}

func updateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResponse{Success: false, Error: "invalid id"})
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, ApiResponse{Success: false, Error: "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}

	var updateData Product
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, ApiResponse{Success: false, Error: "invalid request body"})
	}

	if err := db.Model(&product).Updates(updateData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, ApiResponse{Success: true, Data: product})
}

func deleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResponse{Success: false, Error: "invalid id"})
	}

	var product Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, ApiResponse{Success: false, Error: "product not found"})
		}
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}

	if err := db.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, ApiResponse{Success: false, Error: err.Error()})
	}

	return c.JSON(http.StatusOK, ApiResponse{Success: true, Data: product})
}

func deleteProducts(listProducts []uint) error {
	if listProducts == nil {
		return nil
	}
	if err := db.Where("id IN ?", listProducts).Delete(&Product{}).Error; err != nil {
		return err
	}
	return nil
}
