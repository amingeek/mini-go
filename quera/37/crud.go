package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Color       string    `json:"color"`
	Category    string    `json:"category"`
	CreatedDate time.Time `json:"created_date"`
}

var (
	db = initDB()
)

func createProduct(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return err
	}

	validations := map[string]bool{
		"name is required":         product.Name == "",
		"price is required":        product.Price == 0,
		"color is required":        product.Color == "",
		"category is required":     product.Category == "",
		"created_date is required": product.CreatedDate.IsZero(),
	}

	for msg, failed := range validations {
		if failed {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": msg})
		}
	}

	var existing Product
	if err := db.Where("name = ?", product.Name).First(&existing).Error; err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "product with this name already exists",
		})
	}

	if err := db.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, product)
}

func getProducts(c echo.Context) error {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}
func getCategory(c echo.Context) error {
	var products []Product
	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "category is required"})
	}
	if err := db.Where("category = ?", category).Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	var product Product
	if c.Param("id") == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}
	id := c.Param("id")

	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, product)
}

func updateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "product not found"})
	}

	var updateData Product
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := db.Model(&product).Updates(updateData).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}
