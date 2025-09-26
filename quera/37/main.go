package main

import (
	"net/http"
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

	if product.Name == "" {
		return c.JSON(http.StatusBadRequest, "name is required")
	}
	if product.Price == 0 {
		return c.JSON(http.StatusBadRequest, "price is required")
	}
	if product.Color == "" {
		return c.JSON(http.StatusBadRequest, "color is required")
	}
	if product.Category == "" {
		return c.JSON(http.StatusBadRequest, "category is required")
	}
	if product.CreatedDate.IsZero() {
		return c.JSON(http.StatusBadRequest, "created_date is required")
	}

	if err := db.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, product)
}

func main() {
	e := echo.New()
	e.POST("/products", createProduct)
	//e.GET("/products", getProducts)
	//e.GET("/products/:id", getProduct)
	//e.PUT("/products/:id", updateProduct)
	//e.DELETE("/products/:id", deleteProduct)
	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
