package zokyo

import (
	"MusicDev33/mdapi3/internal/database"
	"MusicDev33/mdapi3/internal/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetConversationsRoute retrieves all conversations for a user
func GetConversationsRoute(c *gin.Context) {
	username := c.Param("username")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.DB.Collection("conversations")

	// Find all conversations for the user, sorted by _id descending, limit 500
	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}).SetLimit(500)
	cursor, err := collection.Find(ctx, bson.M{"user": username}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something went wrong"})
		return
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": conversations})
}

// GetChatsByConvIdRoute retrieves all chats for a conversation
func GetChatsByConvIdRoute(c *gin.Context) {
	convId := c.Param("convId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.DB.Collection("chats")

	// Find all chats for the conversation, sorted by timestamp ascending, limit 500
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}).SetLimit(500)
	cursor, err := collection.Find(ctx, bson.M{"conversationId": convId}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something went wrong"})
		return
	}
	defer cursor.Close(ctx)

	var chats []models.Chat
	if err = cursor.All(ctx, &chats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": chats})
}
