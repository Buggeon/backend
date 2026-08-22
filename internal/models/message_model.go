package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	SenderID  primitive.ObjectID `bson:"sender_id" json:"sender_id"`
	CardID    primitive.ObjectID `bson:"card_id" json:"card_id"`
	ReplyTo   primitive.ObjectID `bson:"reply_to" json:"reply_to"`
	Content   string             `bson:"content" json:"content"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
