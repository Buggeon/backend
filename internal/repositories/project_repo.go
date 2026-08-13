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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectRepo struct {
	collection *mongo.Collection
}

func NewProjectRepo() *ProjectRepo {
	return &ProjectRepo{
		collection: db.GetCollection("projects"),
	}
}

func (r *ProjectRepo) CreateProject(project *models.Project) error {

	project.ID = primitive.NewObjectID()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.TODO(), project)
	return err

}

func (r *ProjectRepo) DeleteProject(projectID primitive.ObjectID) error {

	result, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": projectID})

	if result.DeletedCount != 1 || err != nil {
		return errors.New("Failed to delete card")
	}

	return nil

}

func (r *ProjectRepo) GetProject(projectID primitive.ObjectID) (models.Project, error) {

	var project models.Project

	err := r.collection.FindOne(context.TODO(), bson.M{"_id": projectID}).Decode(&project)

	return project, err

}

func (r *ProjectRepo) GetProjects(amount int, userID primitive.ObjectID) ([]models.Project, error) {

	opts := options.Find().SetLimit(int64(amount)).SetSort(bson.D{{"created_at", 1}})

	var projects []models.Project
	filter := bson.M{
		"members": bson.M{"$in": []primitive.ObjectID{userID}},
	}

	cursor, err := r.collection.Find(context.TODO(), filter, opts)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &projects); err != nil {
		return nil, err
	}

	return projects, nil

}
