package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Create Struct
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Password string             `json:"password" bson:"password"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Url       string             `json:"url,omitempty" bson:"url,omitempty"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	User      *User              `json:"user" bson:"user,omitempty"`
}
