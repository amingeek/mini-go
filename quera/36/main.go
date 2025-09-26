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
	CreatedDate time.Time `json:"created_date"`
}

var (
	db = initializeDB()
)

// var products []Product = []Product{}
func createProduct(c echo.Context) error {
	product := new(Product)

	// Bind incoming JSON -> struct
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	// Save product in DB with GORM
	if err := db.Create(product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	// Return the saved product as JSON
	return c.JSON(http.StatusCreated, product)
}

func getProducts(c echo.Context) error {
	products := []*Product{}
	if err := db.Find(&products).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, products)
}

func getProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid product ID",
		})
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "product not found",
		})
	}

	return c.JSON(http.StatusOK, product)
}

func updateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid product ID",
		})
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "product not found",
		})
	}

	var updateData Product
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	product = updateData

	if err := db.Save(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, product)
}

func deleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid product ID",
		})
	}
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "product not found",
		})
	}
	if err := db.Delete(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, product)
}

func main() {
	e := echo.New()
	e.POST("/products", createProduct)
	e.GET("/products", getProducts)
	e.GET("/products/:id", getProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)
	err := e.Start(":8080")
	if err != nil {
		e.Logger.Fatal(err)
	}
}
