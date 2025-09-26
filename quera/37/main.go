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

//func getProduct(c echo.Context) error {
//	var product Product
//	id := c.Param("id")
//
//}

func main() {
	e := echo.New()
	e.POST("/products", createProduct)
	e.GET("/products", getProducts)
	e.GET("/products/category/:category", getCategory)
	//e.GET("/products/:id", getProduct)
	//e.PUT("/products/:id", updateProduct)
	//e.DELETE("/products/:id", deleteProduct)
	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}
}
