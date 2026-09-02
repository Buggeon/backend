package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schema struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	ProjectID primitive.ObjectID `bson:"project_id" json:"project_id"`
	Direction string             `bson:"direction" json:"direction"`
	Url       string             `bson:"url" json:"url"`
	Name      string             `bson:"name" json:"name"`
	AuthorID  primitive.ObjectID `bson:"author_id" json:"authorId"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
