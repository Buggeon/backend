package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Token     string             `bson:"token" json:"token"`
	Type      string             `bson:"type" json:"type"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	Revoked   bool               `bson:"revoked" json:"revoked"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
}
