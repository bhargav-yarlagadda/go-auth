package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct for MongoDB
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `json:"name" validate:"required"`
	Email       string             `json:"email" validate:"required,email"`
	Password    string             `json:"password" validate:"required"`
	CreatedTime time.Time          `bson:"createdTime,omitempty" json:"-"`
}

type SignInData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}