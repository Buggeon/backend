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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		collection: db.GetCollection("users"),
	}
}

func (r *UserRepo) Create(user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(context.TODO(), user)
	return err
}

func (r *UserRepo) GetByLogin(userLogin string) (*models.User, error) {

	var result models.User

	err := r.collection.FindOne(context.TODO(), bson.M{"login": userLogin}).Decode(&result)

	return &result, err
}

func (r *UserRepo) GetUserPasswordHash(userLogin string) (string, error) {

	var result struct {
		Password string `bson:"password"`
	}

	err := r.collection.FindOne(context.TODO(), bson.M{"login": userLogin}).Decode(&result)

	return result.Password, err

}

func (r *UserRepo) GetAllUsers() ([]models.User, error) {

	var users []models.User

	cursor, err := r.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil

}
