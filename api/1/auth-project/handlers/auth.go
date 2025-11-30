package handlers

import (
	"api/database"
	"api/models"
	"api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}

	hashedPassword, err := utils.PasswordHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	userResponse := models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Register Successfully", userResponse)
}

func Login(c *gin.Context) {
	var req models.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		return
	}

	if err := utils.PasswordCompare(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password incorrect"})
		return
	}

	userResponse := models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	utils.SuccessResponse(c, http.StatusOK, "Login Successfully", userResponse)
}
