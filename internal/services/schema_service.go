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

package services

import (
	s3storage "bugtracker/internal/S3Storage"
	"bugtracker/internal/models"
	"bugtracker/internal/repositories"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/loveholidays/excalidraw-decrypt/pkg/excalidrawdecrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SchemaService struct {
	projectRepo *repositories.ProjectRepo
	schemaRepo  *repositories.SchemaRepo
	s3Storage   *s3storage.S3Storage
}

func NewSchemaService(
	projectRepo *repositories.ProjectRepo,
	schemaRepo *repositories.SchemaRepo,
	s3Storage *s3storage.S3Storage,
) *SchemaService {
	return &SchemaService{
		projectRepo: projectRepo,
		schemaRepo:  schemaRepo,
		s3Storage:   s3Storage,
	}
}

func (s *SchemaService) CreateSchema(projectID, direction, remoteUrl, name, authorID string) error {

	schemaID := primitive.NewObjectID()

	projectObjID, err := primitive.ObjectIDFromHex(projectID)

	if err != nil {
		return err
	}

	authorObjID, err := primitive.ObjectIDFromHex(authorID)

	if err != nil {
		return err
	}

	_, err = s.schemaRepo.CreateSchema(&models.Schema{
		ID:        schemaID,
		ProjectID: projectObjID,
		Direction: direction,
		Url:       remoteUrl,
		Name:      name,
		AuthorID:  authorObjID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	shareableID := strings.Split(remoteUrl, "json=")[1]

	decrypter := excalidrawdecrypt.CreateShareableExcalidrawDecrypter()

	plaintext, err := decrypter.Decrypt(shareableID)
	schemaUrl, err := s.s3Storage.Upload(context.Background(), fmt.Sprintf("projects/%s/schemas/%s", projectID, schemaID.Hex()), strings.NewReader(plaintext), "application/json")

	if err != nil {
		return err
	}

	err = s.projectRepo.AddSchema(projectObjID, schemaUrl)

	if err != nil {
		return err
	}

	return nil

}

func (s *SchemaService) GetSchemas(projectID string) ([]models.Schema, error) {

	projectObjID, err := primitive.ObjectIDFromHex(projectID)

	if err != nil {
		return []models.Schema{}, err
	}

	return s.schemaRepo.GetSchemas(projectObjID)

}

func (s *SchemaService) GetSchema(schemaID string) (models.Schema, error) {

	schemaObjID, err := primitive.ObjectIDFromHex(schemaID)

	if err != nil {
		return models.Schema{}, err
	}

	return s.schemaRepo.GetSchema(schemaObjID)

}
