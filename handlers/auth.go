package handlers

import (
	"log"
	"net/http"
	"todo-list/config"
	"todo-list/models"
	"todo-list/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles POST /register
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// 1. Check for unique email
	var existingUser models.User
	if config.DB.Where("email = ?", user.Email).First(&existingUser).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered"})
		return
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	log.Printf("Login Request: Email=%s, Password=%s", user.Email, user.Password)


	// 3. Save User
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// 4. Generate Token and respond
	token, _ := utils.GenerateToken(user.ID)
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

// LoginUser handles POST /login
func LoginUser(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	log.Printf("Login Request: Email=%s, Password=%s", loginRequest.Email, loginRequest.Password)

	// 1. Find User by Email
	var user models.User
	if err := config.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found with email: " + loginRequest.Email,})
		return
	}
	log.Printf("User retrieved: ID=%d, Hashed Password=%s", user.ID, user.Password)
	log.Printf("Password from request: %s", loginRequest.Password)

	// 2. Verify Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Password Incorrect, err : " + err.Error()})
		return
	}

	// 3. Generate Token and respond
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}