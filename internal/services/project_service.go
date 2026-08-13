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

type ProjectService struct {
	projectRepo   *repositories.ProjectRepo
	memberService *MemberService
}

func NewProjectService(projectRepo *repositories.ProjectRepo, memberService *MemberService) *ProjectService {
	return &ProjectService{
		projectRepo:   projectRepo,
		memberService: memberService,
	}
}

func (s *ProjectService) CreateProject(projectData *dto.CreateProjectDto) error {

	projectID := primitive.NewObjectID()
	ownerId, err := primitive.ObjectIDFromHex(projectData.UserId)

	if err != nil {
		return err
	}

	project := &models.Project{
		Name:        projectData.Name,
		Description: projectData.Description,
		OwnerId:     ownerId,
	}

	if err := s.projectRepo.CreateProject(project); err != nil {
		return err
	}

	for _, member := range projectData.Members {
		if err := s.memberService.CreateMember(&dto.CreateMemberDto{
			UserId:     member.UserId,
			ProjectId:  projectID.Hex(),
			Role:       member.Role,
			Directions: member.Directions,
		}); err != nil {
			return err
		}
	}

	return nil

}

func (s *ProjectService) GetProjects(dto *dto.GetProjectsDto, userID string) ([]models.Project, error) {

	objID, err := primitive.ObjectIDFromHex(userID)

	projects, err := s.projectRepo.GetProjects(dto.Amount, objID)

	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	return projects, nil

}

func (s *ProjectService) GetProject(dto *dto.GetProjectDto) (models.Project, error) {

	objID, err := primitive.ObjectIDFromHex(dto.ProjectID)

	project, err := s.projectRepo.GetProject(objID)

	fmt.Println(err)

	if err != nil {
		return models.Project{}, err
	}

	return project, nil

}

func (s *ProjectService) DeleteProject(dto *dto.GetProjectDto) error {

	objID, err := primitive.ObjectIDFromHex(dto.ProjectID)

	err = s.projectRepo.DeleteProject(objID)

	return err

}
