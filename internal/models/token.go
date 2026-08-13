package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Token     string             `bson:"token" json:"token"`
	Type      string             `bson:"type" json:"type"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Revoked   bool               `bson:"revoked" json:"revoked"`
}

type TokenReponse struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresIn    int
}
