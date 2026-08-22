// Buggeon - SelfHosted service for bug and task tracking
// Copyright (C) 2026 DEVE corp.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package repositories

import (
	"bugtracker/internal/db"
	"bugtracker/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepo struct {
	collection *mongo.Collection
}

func NewMessageRepo() *MessageRepo {
	return &MessageRepo{
		collection: db.GetCollection("messages"),
	}
}

func (r *MessageRepo) NewMessage(message *models.Message) error {

	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(context.TODO(), message)
	return err

}

func (r *MessageRepo) DeleteMessage(messageID primitive.ObjectID) error {

	result, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": messageID})

	if result.DeletedCount != 1 || err != nil {
		return errors.New("Failed to delete message")
	}

	return nil

}

func (r *MessageRepo) GetMessage(messageID primitive.ObjectID) (models.Member, error) {

	var message models.Member

	err := r.collection.FindOne(context.TODO(), bson.M{"_id": messageID}).Decode(&message)

	return message, err

}

func (r *MessageRepo) GetMessages(projectID primitive.ObjectID) ([]models.Member, error) {

	var members []models.Member

	cursor, err := r.collection.Find(context.TODO(), bson.M{"project_id": projectID})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &members); err != nil {
		return nil, err
	}

	return members, nil

}
