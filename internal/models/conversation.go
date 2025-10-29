package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	User      string             `bson:"user" json:"user" binding:"required"`
	Name      string             `bson:"name" json:"name" binding:"required"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
