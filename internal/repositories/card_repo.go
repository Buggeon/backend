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

type CardRepo struct {
	collection *mongo.Collection
}

func NewCardRepo() *CardRepo {
	return &CardRepo{
		collection: db.GetCollection("cards"),
	}
}

func (r *CardRepo) CreateCard(card *models.Card) (primitive.ObjectID, error) {

	card.ID = primitive.NewObjectID()
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.TODO(), card)
	return card.ID, err

}

func (r *CardRepo) DeleteCard(cardID primitive.ObjectID) error {

	result, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": cardID})

	if result.DeletedCount != 1 || err != nil {
		return errors.New("Failed to delete card")
	}

	return nil

}

func (r *CardRepo) GetCard(boardID primitive.ObjectID) (models.Card, error) {

	var card models.Card

	err := r.collection.FindOne(context.TODO(), bson.M{"_id": boardID}).Decode(&card)

	return card, err

}

func (r *CardRepo) GetCards(boardID primitive.ObjectID) ([]models.Card, error) {

	var cards []models.Card

	cursor, err := r.collection.Find(context.TODO(), bson.M{"board_id": boardID})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &cards); err != nil {
		return nil, err
	}

	return cards, nil

}

func (r *CardRepo) AddMessage(cardID, messageID primitive.ObjectID) error {

	filter := bson.M{"_id": cardID}
	update := bson.M{"$push": bson.M{"messages": messageID}}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)

	return err

}
