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

type MemberRepo struct {
	collection *mongo.Collection
}

func NewMemberRepo() *MemberRepo {
	return &MemberRepo{
		collection: db.GetCollection("members"),
	}
}

func (r *MemberRepo) CreateMember(member *models.Member) (primitive.ObjectID, error) {

	member.ID = primitive.NewObjectID()
	member.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(context.TODO(), member)
	return member.ID, err

}

func (r *MemberRepo) DeleteMember(memberID primitive.ObjectID) error {

	result, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": memberID})

	if result.DeletedCount != 1 || err != nil {
		return errors.New("Failed to delete card")
	}

	return nil

}

func (r *MemberRepo) GetMember(memberID primitive.ObjectID) (models.Member, error) {

	var member models.Member

	err := r.collection.FindOne(context.TODO(), bson.M{"_id": memberID}).Decode(&member)

	return member, err

}

func (r *MemberRepo) GetMembers(projectID primitive.ObjectID) ([]models.Member, error) {

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

func (r *MemberRepo) GetProjectsIDsByUser(userID primitive.ObjectID) ([]primitive.ObjectID, error) {

	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var members []models.Member

	if err := cursor.All(context.TODO(), &members); err != nil {
		return nil, err
	}

	projectsIDs := make([]primitive.ObjectID, len(members))

	for i, m := range members {
		projectsIDs[i] = m.ProjectID
	}

	return projectsIDs, nil
}
