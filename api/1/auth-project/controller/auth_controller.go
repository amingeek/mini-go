package controller

import (
	"auth-project/data/request"
	"auth-project/helper"
	"auth-project/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AuthController struct {
	Db       *gorm.DB
	Validate *validator.Validate
}

func NewAuthController(db *gorm.DB, validate *validator.Validate) *AuthController {
	return &AuthController{Db: db, Validate: validate}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var reqBody request.RegisterRequest

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := fmt.Sprintf("Validation failed for field: %s",
			validationErrors[0].Field())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	var existingUser model.User
	result := c.Db.Where("email = ?", reqBody.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := helper.EncryptPassword(reqBody.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to encrypt password"})
		return
	}

	newUser := model.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: hashedPassword,
	}

	if err := c.Db.Create(&newUser).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var reqBody request.LoginRequest

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Validate.Struct(reqBody); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessage := fmt.Sprintf("Validation failed for field: %s",
			validationErrors[0].Field())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	var user model.User
	result := c.Db.Where("email = ?", reqBody.Email).First(&user)
	if result.RowsAffected < 1 {
		ctx.JSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password"})
		return
	}

	if !helper.ComparePassword(reqBody.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized,
			gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := helper.CreateToken(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
