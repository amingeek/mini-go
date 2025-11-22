package main

import (
	"auth-project/config"
	"auth-project/controller"
	"auth-project/model"
	"auth-project/router"
	"log"

	"github.com/go-playground/validator/v10"
)

func main() {
	// اتصال به دیتابیس
	db := config.DatabaseConnection()

	// مایگریشن جدول
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migration completed!")

	// Validator
	validate := validator.New()

	// Controller
	authController := controller.NewAuthController(db, validate)

	// Setup router
	r := router.SetupRouter(authController)

	// اجرای سرور
	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
