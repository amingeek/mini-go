package main

import "github.com/labstack/echo/v4"

func main() {
	e := echo.New()
	e.POST("/products/", createProduct)
	e.GET("/products/", getProducts)
	e.GET("/products/:id", getProduct)
	e.PUT("/products/:id", updateProduct)
	e.DELETE("/products/:id", deleteProduct)
	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}

}
