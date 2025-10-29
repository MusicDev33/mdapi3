package zokyo

import (
	"MusicDev33/mdapi3/internal/database"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteConversationByIdRoute deletes a conversation and all its chats
func DeleteConversationByIdRoute(c *gin.Context) {
	convId := c.Param("convId")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete all chats associated with this conversation
	chatsCollection := database.DB.DB.Collection("chats")
	_, err := chatsCollection.DeleteMany(ctx, bson.M{"conversationId": convId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something went wrong with deleting the chats"})
		return
	}

	// Convert convId string to ObjectID
	objectId, err := primitive.ObjectIDFromHex(convId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "Invalid conversation ID"})
		return
	}

	// Delete the conversation
	convsCollection := database.DB.DB.Collection("conversations")
	_, err = convsCollection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something went wrong with deleting the conversation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "Successfully deleted conversation!"})
}
