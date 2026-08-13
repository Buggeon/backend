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

package handlers

import (
	"bugtracker/internal/dto"
	"bugtracker/internal/services"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *services.ProjectService
	boardService   *services.BoardService
	cardService    *services.CardService
	memberService  *services.MemberService
}

func NewProjectHandler(projectService *services.ProjectService, cardService *services.CardService, boardServices *services.BoardService, memberService *services.MemberService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		boardService:   boardServices,
		cardService:    cardService,
		memberService:  memberService,
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {

	var dto dto.CreateProjectDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(403)
	}

	h.projectService.CreateProject(&dto)
}

func (h *ProjectHandler) EditProject(c *gin.Context) {

}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {

}

func (h *ProjectHandler) GetProject(c *gin.Context) {

}

func (h *ProjectHandler) GetProjects(c *gin.Context) {

	var dto dto.GetProjectsDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(403)
	}

	//userID, _ := c.Get("userID")
	userID := "6a79b2f98860397db9e68552"

	projects, err := h.projectService.GetProjects(&dto, userID)

	if err != nil {
		c.Status(403)

	} else {
		c.JSON(200, projects)
	}

}

func (h *ProjectHandler) CreateBoard(c *gin.Context) {

}

func (h *ProjectHandler) EditBoard(c *gin.Context) {

}

func (h *ProjectHandler) DeleteBoard(c *gin.Context) {

}

func (h *ProjectHandler) GetBoard(c *gin.Context) {

}

func (h *ProjectHandler) GetBoards(c *gin.Context) {

}

func (h *ProjectHandler) CreateCard(c *gin.Context) {

}

func (h *ProjectHandler) EditCard(c *gin.Context) {

}

func (h *ProjectHandler) DeleteCard(c *gin.Context) {

}

func (h *ProjectHandler) GetCard(c *gin.Context) {

}

func (h *ProjectHandler) GetCards(c *gin.Context) {

}

func (h *ProjectHandler) CreateMember(c *gin.Context) {

}

func (h *ProjectHandler) EditMember(c *gin.Context) {

}

func (h *ProjectHandler) DeleteMember(c *gin.Context) {

}

func (h *ProjectHandler) GetMember(c *gin.Context) {

}

func (h *ProjectHandler) GetMembers(c *gin.Context) {

}
