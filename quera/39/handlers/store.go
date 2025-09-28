package handlers

import (
	"TestEcho/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type tmp struct {
	Name  string
	Price int64
	Count int64
}

func ListProducts(c echo.Context, dbConn *gorm.DB) error {
	products := []tmp{}
	dbConn.Table("products").Select("name", "price", "count").Find(&products)
	return c.JSON(http.StatusOK, products)
}

func AddProduct(c echo.Context, dbConn *gorm.DB) error {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Invalid JSON")
	}
	var p models.Product
	p.Name = jsonBody["name"].(string)
	p.Count = int64(jsonBody["count"].(float64))
	p.Price = int64(jsonBody["price"].(float64))

	createdUser := dbConn.Create(&p)
	if createdUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, "Product Creation Failed")
	}
	return c.JSON(http.StatusCreated, p)
}
