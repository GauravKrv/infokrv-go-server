package model

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name      string            `bson:"name" json:"name" validate:"required,min=2,max=100"`
    Email     string            `bson:"email" json:"email" validate:"required,email"`
    Password  string            `bson:"password,omitempty" json:"password,omitempty" validate:"required,min=6"`
    CreatedAt time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}
