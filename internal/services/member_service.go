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
	memberRepo *repositories.MemberRepo
}

func NewMemberService(memberRepo *repositories.MemberRepo) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
	}
}

func (s *MemberService) CreateMember(dto *dto.CreateMemberDto) error {

	user_id, err := primitive.ObjectIDFromHex(dto.UserId)

	if err != nil {
		return err
	}

	member := &models.Member{
		ID:         primitive.NewObjectID(),
		UserId:     user_id,
		Role:       dto.Role,
		Directions: dto.Directions,
		CreatedAt:  time.Now(),
	}

	err = s.memberRepo.CreateMember(member)
	return err

}

func (s *MemberService) GetMember(dto *dto.GetMemberDto) (models.Member, error) {

	objID, err := primitive.ObjectIDFromHex(dto.MemberID)

	if err != nil {
		return models.Member{}, err
	}

	member, err := s.memberRepo.GetMember(objID)

	return member, err

}

func (s *MemberService) DeleteMember(dto *dto.DeleteMemberDto) error {

	objID, err := primitive.ObjectIDFromHex(dto.MemberID)

	if err != nil {
		return err
	}

	err = s.memberRepo.DeleteMember(objID)

	return err

}

func (s *MemberService) GetMembers(dto *dto.GetMembersDto) ([]models.Member, error) {

	objID, err := primitive.ObjectIDFromHex(dto.ProjectID)

	members, err := s.memberRepo.GetMembers(objID)

	return members, err

}
