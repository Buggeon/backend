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

type BoardRepo struct {
	collection *mongo.Collection
}

func NewBoardRepo() *BoardRepo {
	return &BoardRepo{
		collection: db.GetCollection("boards"),
	}
}

func (r *BoardRepo) CreateBoard(board *models.Board) (primitive.ObjectID, error) {

	board.ID = primitive.NewObjectID()
	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.TODO(), board)
	return board.ID, err

}

func (r *BoardRepo) DeleteBoard(boardID primitive.ObjectID) error {

	result, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": boardID})

	if result.DeletedCount != 1 || err != nil {
		return errors.New("Failed to delete board")
	}

	return nil

}

func (r *BoardRepo) GetBoard(boardID primitive.ObjectID) (models.Board, error) {

	var board models.Board

	err := r.collection.FindOne(context.TODO(), bson.M{"_id": boardID}).Decode(&board)

	return board, err

}

func (r *BoardRepo) GetBoards(projectID primitive.ObjectID) ([]models.Board, error) {

	var boards []models.Board

	cursor, err := r.collection.Find(context.TODO(), bson.M{"project_id": projectID})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &boards); err != nil {
		return nil, err
	}

	return boards, nil

}

func (r *BoardRepo) AddCard(boardID, cardID primitive.ObjectID) error {

	filter := bson.M{"_id": boardID}
	update := bson.M{"$push": bson.M{"cards": cardID}}

	_, err := r.collection.UpdateOne(context.TODO(), filter, update)

	return err

}
