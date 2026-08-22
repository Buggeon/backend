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
	memberRepo    *repositories.MemberRepo
	memberService *MemberService
}

func NewProjectService(projectRepo *repositories.ProjectRepo, memberService *MemberService, memberRepo *repositories.MemberRepo) *ProjectService {
	return &ProjectService{
		projectRepo:   projectRepo,
		memberService: memberService,
		memberRepo:    memberRepo,
	}
}

func (s *ProjectService) CreateProject(projectData *dto.CreateProjectDto) error {

	projectID := primitive.NewObjectID()

	leadObjID, err := s.memberService.CreateMember(&dto.CreateMemberDto{
		ProjectID:  projectID.Hex(),
		UserID:     projectData.LeadID,
		Role:       "Lead",
		Directions: []string{"managment"},
	})

	if err != nil {
		return err
	}

	fmt.Println("===LEAD ID FROM PROJECT SERVICE===")
	fmt.Println(leadObjID)

	project := &models.Project{
		ID:          projectID,
		Name:        projectData.Name,
		Description: projectData.Description,
		LeadID:      leadObjID,
	}

	if err := s.projectRepo.CreateProject(project); err != nil {
		return err
	}

	for _, member := range projectData.Members {
		if _, err := s.memberService.CreateMember(&dto.CreateMemberDto{
			UserID:     member.UserID,
			ProjectID:  projectID.Hex(),
			Role:       member.Role,
			Directions: member.Directions,
		}); err != nil {
			return err
		}
	}

	return nil

}

func (s *ProjectService) GetProjects(userID string) ([]models.Project, error) {

	objID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, err
	}

	projectsIDs, err := s.memberRepo.GetProjectsIDsByUser(objID)

	if err != nil {
		return nil, err
	}

	return s.projectRepo.GetProjectsByIDs(projectsIDs)
}

func (s *ProjectService) GetProject(projectID string) (models.Project, error) {

	objID, err := primitive.ObjectIDFromHex(projectID)

	if err != nil {
		return models.Project{}, err
	}

	project, err := s.projectRepo.GetProject(objID)

	if err != nil {
		return models.Project{}, err
	}

	return project, nil

}

func (s *ProjectService) DeleteProject(projectID string) error {

	objID, err := primitive.ObjectIDFromHex(projectID)

	if err != nil {
		return err
	}

	err = s.projectRepo.DeleteProject(objID)

	if err != nil {
		return err
	}

	members, err := s.memberRepo.GetMembers(objID)

	for _, member := range members {
		err := s.memberRepo.DeleteMember(member.ID)

		if err != nil {
			fmt.Println("Failed to delete member")
		}
	}

	return err

}
