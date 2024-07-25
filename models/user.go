package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name      string             `json:"name" validate:"required"`
    Email     string             `json:"email" validate:"required,email"`
    Password  string             `json:"password" validate:"required"`
}

type Credentials struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

