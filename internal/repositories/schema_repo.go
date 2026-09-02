package repositories

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

import (
	"bugtracker/internal/db"
	"bugtracker/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchemaRepo struct {
	collection *mongo.Collection
}

func NewSchemaRepo() *SchemaRepo {
	return &SchemaRepo{
		collection: db.GetCollection("schemas"),
	}
}

func (r *SchemaRepo) CreateSchema(schema *models.Schema) (primitive.ObjectID, error) {

	_, err := r.collection.InsertOne(context.Background(), schema)

	return schema.ID, err

}

func (r *SchemaRepo) GetSchema(schemaID primitive.ObjectID) (models.Schema, error) {

	var schema models.Schema

	err := r.collection.FindOne(context.Background(), bson.M{"_id": schemaID}).Decode(&schema)

	return schema, err

}

func (r *SchemaRepo) GetSchemas(projectID primitive.ObjectID) ([]models.Schema, error) {

	var schemas []models.Schema

	cursor, err := r.collection.Find(context.Background(), bson.M{"project_id": projectID})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &schemas); err != nil {
		return nil, err
	}

	return schemas, nil

}
