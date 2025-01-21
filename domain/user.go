package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	CollectionUser      = "users"
	AdminRole      Role = "ADMIN"
	UserRole       Role = "USER"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Full_Name       string             `bson:"full_name" json:"full_name"`
	Email           string             `bson:"email" json:"email"`
	Password        string             `bson:"password" json:"password" `
	Role            string             `bson:"role" json:"role"`
	Bio             string             `bson:"bio" json:"bio"`
	To_whom         string             `bson:"to_whom" json:"to_whom"`
	Verify          bool               `bson:"verify" json:"verify"`
	Profile_Picture string             `bson:"profile_picture" 
	json:"profile_picture"`
	Created_At      primitive.DateTime `bson:"created_at" json:"created_at"`
}
