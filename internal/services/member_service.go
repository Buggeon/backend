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

type MemberService struct {
	memberRepo  *repositories.MemberRepo
	projectRepo *repositories.ProjectRepo
}

func NewMemberService(memberRepo *repositories.MemberRepo, projectRepo *repositories.ProjectRepo) *MemberService {
	return &MemberService{
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
	}
}

func (s *MemberService) CreateMember(dto *dto.CreateMemberDto) (primitive.ObjectID, error) {

	userID, err := primitive.ObjectIDFromHex(dto.UserID)

	if err != nil {
		return primitive.ObjectID{}, err
	}

	projectID, err := primitive.ObjectIDFromHex(dto.ProjectID)

	if err != nil {
		return primitive.ObjectID{}, err
	}

	member := &models.Member{
		UserID:     userID,
		Role:       dto.Role,
		ProjectID:  projectID,
		Directions: dto.Directions,
		CreatedAt:  time.Now(),
	}

	memberID, err := s.memberRepo.CreateMember(member)

	return memberID, s.projectRepo.AddMember(projectID, memberID)
}

func (s *MemberService) GetMember(memberID string) (models.Member, error) {

	objID, err := primitive.ObjectIDFromHex(memberID)

	if err != nil {
		return models.Member{}, err
	}

	member, err := s.memberRepo.GetMember(objID)

	return member, err

}

func (s *MemberService) DeleteMember(memberID string) error {

	objID, err := primitive.ObjectIDFromHex(memberID)

	if err != nil {
		return err
	}

	err = s.memberRepo.DeleteMember(objID)

	return err

}

func (s *MemberService) GetMembers(projectID string) ([]models.Member, error) {

	objID, err := primitive.ObjectIDFromHex(projectID)

	if err != nil {
		return nil, err
	}

	members, err := s.memberRepo.GetMembers(objID)

	return members, err

}
