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
	"bugtracker/internal/dto"
	"bugtracker/internal/models"
	"bugtracker/internal/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardService struct {
	cardRepo *repositories.CardRepo
}

func NewCardService(cardRepo *repositories.CardRepo) *CardService {
	return &CardService{
		cardRepo: cardRepo,
	}
}

func (b *CardService) CreateCard(dto *dto.CreateCardDto) error {

	boardID, err := primitive.ObjectIDFromHex(dto.BoardID)
	var assigneeIDs []primitive.ObjectID

	for _, assigneeID := range dto.Assignees {
		objID, err := primitive.ObjectIDFromHex(assigneeID)

		if err != nil {
			continue
		}

		assigneeIDs = append(assigneeIDs, objID)
	}

	if err != nil {
		return err
	}

	card := &models.Card{
		Title:       dto.Title,
		Description: dto.Description,
		Content:     dto.Content,
		Assignees:   assigneeIDs,
		BoardID:     boardID,
		Priority:    dto.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = b.cardRepo.CreateCard(card)

	return err
}

func (b *CardService) GetCard(boardID string) (models.Card, error) {

	objID, err := primitive.ObjectIDFromHex(boardID)

	if err != nil {
		return models.Card{}, err
	}

	card, err := b.cardRepo.GetCard(objID)

	return card, err
}

func (b *CardService) GetCards(projectID string) ([]models.Card, error) {

	objID, err := primitive.ObjectIDFromHex(projectID)

	cards, err := b.cardRepo.GetCards(objID)

	return cards, err
}
