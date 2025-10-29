package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ZUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string             `bson:"username" json:"username" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required,min=8"`
	ModelPref string             `bson:"modelPref,omitempty" json:"modelPref,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
