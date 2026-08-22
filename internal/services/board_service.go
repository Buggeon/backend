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
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	boardRepo   *repositories.BoardRepo
	projectRepo *repositories.ProjectRepo
}

func NewBoardService(boardRepo *repositories.BoardRepo, projectRepo *repositories.ProjectRepo) *BoardService {
	return &BoardService{
		boardRepo:   boardRepo,
		projectRepo: projectRepo,
	}
}

func (b *BoardService) CreateBoard(dto *dto.CreateBoardDto) error {

	projectID, err := primitive.ObjectIDFromHex(dto.ProjectID)

	if err != nil {
		return err
	}

	board := &models.Board{
		Name:      dto.Name,
		ProjectID: projectID,
		Direction: dto.Direction,
		Cards:     []primitive.ObjectID{},
	}

	boardID, err := b.boardRepo.CreateBoard(board)

	if err != nil {
		return err
	}

	return b.projectRepo.AddBoard(projectID, boardID)

}

func (b *BoardService) GetBoard(boardID string) (models.Board, error) {

	objID, err := primitive.ObjectIDFromHex(boardID)

	if err != nil {
		return models.Board{}, err
	}

	boards, err := b.boardRepo.GetBoard(objID)

	return boards, err
}

func (b *BoardService) GetBoards(projectID string) ([]models.Board, error) {

	objID, err := primitive.ObjectIDFromHex(projectID)

	boards, err := b.boardRepo.GetBoards(objID)

	fmt.Println(err)

	return boards, err
}

func (b *BoardService) DeleteBoard(boardID string) error {

	objID, err := primitive.ObjectIDFromHex(boardID)

	if err != nil {
		return err
	}

	err = b.boardRepo.DeleteBoard(objID)

	return err
}
