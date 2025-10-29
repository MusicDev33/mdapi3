package zokyo

import (
	"MusicDev33/mdapi3/internal/config"
	"MusicDev33/mdapi3/internal/database"
	"MusicDev33/mdapi3/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserData struct {
	Username string             `json:"username"`
	ID       primitive.ObjectID `json:"_id"`
}

// AuthRoute handles user authentication
func AuthRoute(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user by username
	var user models.ZUser
	collection := database.DB.DB.Collection("zusers")
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": GetUserNotFoundMsg()})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": GetPassMsg()})
		return
	}

	// Return user data without password
	userData := UserData{
		Username: user.Username,
		ID:       user.ID,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": userData})
}

// CreateLoginRoute creates a new user login
func CreateLoginRoute(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false})
		return
	}

	cfg := config.Get()

	// Check if username is in whitelist
	isWhitelisted := false
	for _, u := range cfg.WhitelistUsers {
		if u == req.Username {
			isWhitelisted = true
			break
		}
	}

	if !isWhitelisted {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	// Create new user
	newUser := models.ZUser{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := database.DB.DB.Collection("zusers")
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false})
		return
	}

	userData := UserData{
		Username: newUser.Username,
		ID:       newUser.ID,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": userData})
}

// CheckWhitelistUserRoute checks if a username is available and whitelisted
func CheckWhitelistUserRoute(c *gin.Context) {
	username := c.Param("username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if user already exists
	var user models.ZUser
	collection := database.DB.DB.Collection("zusers")
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == nil {
		// User already exists
		c.JSON(http.StatusOK, gin.H{"success": false})
		return
	}

	cfg := config.Get()

	// Check if username is in whitelist
	for _, u := range cfg.WhitelistUsers {
		if u == username {
			c.JSON(http.StatusOK, gin.H{"success": true})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": false})
}
