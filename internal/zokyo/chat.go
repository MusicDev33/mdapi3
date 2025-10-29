package zokyo

import (
	"MusicDev33/mdapi3/internal/database"
	"MusicDev33/mdapi3/internal/models"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tokenThreshold = 4096 * 2

type CreateChatRequest struct {
	User   string     `json:"user" binding:"required"`
	ConvID string     `json:"convId"`
	Msg    string     `json:"msg" binding:"required"`
	Mode   string     `json:"mode" binding:"required"`
	Engine ChatEngine `json:"engine"`
}

// HandleSecureEngine enforces security constraints on engine selection
func HandleSecureEngine(ip string, engine ChatEngine) ChatEngine {
	// If IP requires more security, disable DeepSeek
	if engine == EngineDeepSeek {
		return EngineClaude
	}
	return engine
}

// CreateNewChatRoute handles the creation of new chat messages and AI responses
func CreateNewChatRoute(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "Invalid request"})
		return
	}

	// Validate message content
	if strings.TrimSpace(req.Msg) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "Message cannot be empty"})
		return
	}

	// Default engine to Claude if not specified or invalid
	if req.Engine == "" {
		req.Engine = EngineClaude
	}

	// Apply security constraints
	req.Engine = HandleSecureEngine(c.ClientIP(), req.Engine)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var conv *models.Conversation
	convID := req.ConvID

	// Create new conversation if convId is empty
	if convID == "" {
		convName := GenerateName()
		newConv := models.Conversation{
			ID:        primitive.NewObjectID(),
			User:      req.User,
			Name:      convName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		convsCollection := database.DB.DB.Collection("conversations")
		result, err := convsCollection.InsertOne(ctx, newConv)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something broke with Zokyo's backend"})
			return
		}

		newConv.ID = result.InsertedID.(primitive.ObjectID)
		conv = &newConv
		convID = newConv.ID.Hex()
	}

	// Set agent configuration based on mode
	temperature := 1.0
	topP := 1.0

	var systemMessages []Message

	if req.Mode == "code" {
		temperature = 0.3
		systemMessages = []Message{
			{Role: "user", Content: "You are a terse code completion machine. You will answer future questions with just code and nothing more. Do not bother explaining what the code does."},
			{Role: "assistant", Content: "Okay, I will answer your future questions with just code snippets, and no other text. Let's get started!"},
		}
	}

	// Fetch previous chats
	chatsCollection := database.DB.DB.Collection("chats")
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}}).SetLimit(500)
	cursor, err := chatsCollection.Find(ctx, bson.M{"conversationId": convID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something broke with Zokyo's backend"})
		return
	}
	defer cursor.Close(ctx)

	var prevChats []models.Chat
	if err = cursor.All(ctx, &prevChats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something broke with Zokyo's backend"})
		return
	}

	// Build message list
	var messages []Message
	messages = append(messages, systemMessages...)

	for _, chat := range prevChats {
		messages = append(messages, Message{
			Role:    chat.Role,
			Content: chat.Content,
		})
	}

	messages = append(messages, Message{
		Role:    "user",
		Content: req.Msg,
	})

	// Ensure we don't exceed token threshold
	for !CountTokens(messages, tokenThreshold) && len(prevChats) > 0 {
		prevChats = prevChats[1:]
		messages = []Message{}
		messages = append(messages, systemMessages...)

		for _, chat := range prevChats {
			messages = append(messages, Message{
				Role:    chat.Role,
				Content: chat.Content,
			})
		}

		messages = append(messages, Message{
			Role:    "user",
			Content: req.Msg,
		})
	}

	// Generate AI response
	agentConfig := AgentConfig{
		TopP:        topP,
		Temperature: temperature,
	}

	content, err := GenerateChat(req.Engine, messages, agentConfig)
	if err != nil {
		fmt.Printf("AI generation error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": "Something broke with Zokyo's backend"})
		return
	}

	// Create chat records
	timestamp := time.Now().UnixMilli()

	userChat := models.Chat{
		ID:             primitive.NewObjectID(),
		ConversationID: convID,
		Role:           "user",
		Content:        req.Msg,
		Timestamp:      timestamp,
	}

	assistantChat := models.Chat{
		ID:             primitive.NewObjectID(),
		ConversationID: convID,
		Role:           "assistant",
		Content:        content,
		Timestamp:      timestamp + 1,
	}

	// Save both chats
	_, err = chatsCollection.InsertOne(ctx, userChat)
	if err != nil {
		fmt.Printf("Error saving user chat: %v\n", err)
	}

	_, err = chatsCollection.InsertOne(ctx, assistantChat)
	if err != nil {
		fmt.Printf("Error saving assistant chat: %v\n", err)
	}

	// Build response
	responseData := gin.H{
		"success": true,
		"msg":     "Successfully received response.",
		"newChat": assistantChat,
	}

	if conv != nil {
		responseData["newConversation"] = conv
	}

	c.JSON(http.StatusOK, responseData)
}
