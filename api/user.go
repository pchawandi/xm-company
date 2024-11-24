package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pchawandi/xm-company/auth"
	"github.com/pchawandi/xm-company/database"
	"github.com/pchawandi/xm-company/models"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
}

type userRepository struct {
	DB  database.Database
	Ctx context.Context
}

func NewUserRepository(ctx context.Context, db database.Database) *userRepository {
	return &userRepository{
		DB:  db,
		Ctx: ctx,
	}
}

func (r *userRepository) LoginHandler(c *gin.Context) {
	var incomingUser models.User
	var dbUser models.User

	// Get JSON body
	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Fetch the user from the database
	if err := r.DB.Where("username = ?", incomingUser.Username).First(&dbUser).Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(dbUser.Username, dbUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRepository) RegisterHandler(c *gin.Context) {
	var user models.LoginUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Create new user
	newUser := models.User{Username: user.Username, Password: hashedPassword, Role: user.Role}

	// Save the user to the database
	if err := r.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not save user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}
