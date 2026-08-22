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
	cardRepo  *repositories.CardRepo
	boardRepo *repositories.BoardRepo
}

func NewCardService(cardRepo *repositories.CardRepo, boardRepo *repositories.BoardRepo) *CardService {
	return &CardService{
		cardRepo:  cardRepo,
		boardRepo: boardRepo,
	}
}

func (c *CardService) CreateCard(dto *dto.CreateCardDto) error {

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
		Title:     dto.Title,
		Content:   dto.Content,
		Assignees: assigneeIDs,
		BoardID:   boardID,
		Priority:  dto.Priority,
		Messages:  []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	cardID, err := c.cardRepo.CreateCard(card)

	return c.boardRepo.AddCard(boardID, cardID)
}

func (c *CardService) GetCard(cardID string) (models.Card, error) {

	objID, err := primitive.ObjectIDFromHex(cardID)

	if err != nil {
		return models.Card{}, err
	}

	card, err := c.cardRepo.GetCard(objID)

	return card, err
}

func (c *CardService) DeleteCard(cardID string) error {

	objID, err := primitive.ObjectIDFromHex(cardID)

	if err != nil {
		return err
	}

	err = c.cardRepo.DeleteCard(objID)

	return err
}

func (c *CardService) GetCards(boardID string) ([]models.Card, error) {

	objID, err := primitive.ObjectIDFromHex(boardID)

	cards, err := c.cardRepo.GetCards(objID)

	return cards, err
}

func (c *CardService) CloseCard() {

}
