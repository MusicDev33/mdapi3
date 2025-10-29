package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ConversationID string             `bson:"conversationId" json:"conversationId" binding:"required"`
	Role           string             `bson:"role" json:"role" binding:"required"`
	Content        string             `bson:"content" json:"content" binding:"required"`
	Timestamp      int64              `bson:"timestamp" json:"timestamp" binding:"required"`
}
