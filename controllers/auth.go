package controllers

import (
	"backend/confiq"
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var JwtKey = []byte("rahasia")

func Register(c *gin.Context) {

	var input dto.RegisterInput
	
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Bind error:", err) 
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var existingUser models.User
	if err := confiq.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Error selain "record not found" (misal error koneksi DB)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encryption failed"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := confiq.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or DB error"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully",
                 "user":dto.UserResponse{
					Name: user.Name,
				 }})
}

func Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := confiq.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenStr, err := token.SignedString(JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr,"user":dto.UserResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
	}})
}