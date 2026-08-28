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
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *services.ProjectService
	boardService   *services.BoardService
	cardService    *services.CardService
	memberService  *services.MemberService
	messageService *services.MessageService
}

func NewProjectHandler(
	projectService *services.ProjectService,
	cardService *services.CardService,
	boardServices *services.BoardService,
	memberService *services.MemberService,
	messageService *services.MessageService,
) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		boardService:   boardServices,
		cardService:    cardService,
		memberService:  memberService,
		messageService: messageService,
	}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {

	leadID := c.Request.FormValue("lead_id")
	name := c.Request.FormValue("name")
	description := c.Request.FormValue("description")
	//members := c.Request.FormValue("members")
	logo, _ := c.FormFile("logo")
	progress := c.Request.FormValue("progress")

	progressInt, err := strconv.Atoi(progress)

	if err != nil {
		c.Status(403)
		c.Abort()
		return
	}

	err = h.projectService.CreateProject(&dto.CreateProjectDto{
		LeadID:      leadID,
		Name:        name,
		Description: description,
		Logo:        logo,
		Progress:    progressInt,
		//Members: string(membersJSON),
	})

	if err != nil {
		c.Status(403)
		c.Abort()
		return
	}

	c.Status(201)

}

func (h *ProjectHandler) SetProjectLogo(c *gin.Context) {

	file, err := c.FormFile("logo")
	projectID := c.Param("project_id")

	if err != nil {
		c.JSON(403, "New logo file was not provide")
		return
	}

	err = h.projectService.SetProjectLogo(projectID, file)

	if err != nil {
		c.JSON(500, "Failed to upload new logo file")
		return
	}

}

func (h *ProjectHandler) SetProjectName(c *gin.Context) {

}

func (h *ProjectHandler) SetProjectDescription(c *gin.Context) {

}

func (h *ProjectHandler) SetProjectProgress(c *gin.Context) {

}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {

	projectID := c.Param("project_id")

	err := h.projectService.DeleteProject(projectID)

	if err != nil {
		c.Status(403)
		c.Abort()
		return
	}

	c.Status(200)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {

	projectID := c.Param("project_id")

	if project, err := h.projectService.GetProject(projectID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, project)
	}
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {

	userID, _ := c.Get("userID")

	projects, err := h.projectService.GetProjects(userID.(string))

	if err != nil {
		fmt.Println(err)
		c.Status(403)
		c.Abort()
		return

	} else {
		c.JSON(200, projects)
	}

}

func (h *ProjectHandler) CreateBoard(c *gin.Context) {

	var dto dto.CreateBoardDto
	projectID := c.Param("project_id")

	if err := c.ShouldBindJSON(&dto); err != nil {
		fmt.Println(err)
		c.Status(403)
		c.Abort()
		return
	}

	dto.ProjectID = projectID

	if err := h.boardService.CreateBoard(&dto); err != nil {

		fmt.Println(err)

		c.Status(403)
		c.Abort()
		return
	} else {
		c.Status(200)
	}

}

func (h *ProjectHandler) EditBoard(c *gin.Context) {

}

func (h *ProjectHandler) DeleteBoard(c *gin.Context) {

	boardID := c.Param("board_id")

	if err := h.boardService.DeleteBoard(boardID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.Status(200)
	}
}

func (h *ProjectHandler) GetBoard(c *gin.Context) {

	boardID := c.Param("board_id")

	if board, err := h.boardService.GetBoard(boardID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, board)
	}
}

func (h *ProjectHandler) GetBoards(c *gin.Context) {

	projectID := c.Param("project_id")

	if boards, err := h.boardService.GetBoards(projectID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, boards)
	}
}

func (h *ProjectHandler) CreateCard(c *gin.Context) {
	var dto dto.CreateCardDto

	boardID := c.Param("board_id")

	if err := c.ShouldBindJSON(&dto); err != nil {
		fmt.Println(err)
		c.Status(403)
		c.Abort()
		return
	}

	dto.BoardID = boardID

	if err := h.cardService.CreateCard(&dto); err != nil {
		fmt.Println(err)
		c.Status(403)
		c.Abort()
		return
	} else {
		c.Status(200)
	}
}

func (h *ProjectHandler) EditCard(c *gin.Context) {
	// var dto dto.CreateCardDto

	// if err := c.ShouldBindJSON(&dto); err != nil {
	// 	c.Status(403)
	// 	c.Abort()
	// }

	// if err := h.boardService.CreateCard(&dto); err != nil {
	// 	c.Status(403)
	// 	c.Abort()
	// } else {
	// 	c.Status(200)
	// }
}

func (h *ProjectHandler) DeleteCard(c *gin.Context) {

	cardID := c.Param("card_id")

	if err := h.cardService.DeleteCard(cardID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.Status(200)
	}
}

func (h *ProjectHandler) GetCard(c *gin.Context) {

	cardID := c.Param("card_id")

	if card, err := h.cardService.GetCard(cardID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, card)
	}
}

func (h *ProjectHandler) GetCards(c *gin.Context) {

	boardID := c.Param("board_id")

	if cards, err := h.cardService.GetCards(boardID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, cards)
	}
}

func (h *ProjectHandler) CreateMember(c *gin.Context) {
	var dto dto.CreateMemberDto

	projectID := c.Param("project_id")

	if err := c.ShouldBindJSON(&dto); err != nil {
		fmt.Println(err)
		c.Status(403)
		c.Abort()
		return
	}

	dto.ProjectID = projectID

	if _, err := h.memberService.CreateMember(&dto); err != nil {
		fmt.Println(err)
		c.Status(403)
		c.Abort()
		return
	} else {
		c.Status(200)
	}
}

func (h *ProjectHandler) EditMember(c *gin.Context) {
	// var dto dto.CreateBoardDto

	// if err := c.ShouldBindJSON(&dto); err != nil {
	// 	c.Status(403)
	// }

	// if err := h.boardService.CreateBoard(&dto); err != nil {
	// 	c.Status(403)
	// } else {
	// 	c.Status(200)
	// }
}

func (h *ProjectHandler) DeleteMember(c *gin.Context) {

	memberID := c.Param("member_id")

	if err := h.memberService.DeleteMember(memberID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.Status(200)
	}
}

func (h *ProjectHandler) GetMember(c *gin.Context) {

	memberID := c.Param("member_id")

	if member, err := h.memberService.GetMember(memberID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, member)
	}
}

func (h *ProjectHandler) GetMembers(c *gin.Context) {

	projectID := c.Param("project_id")

	if members, err := h.memberService.GetMembers(projectID); err != nil {
		c.Status(403)
		c.Abort()
		return
	} else {
		c.JSON(200, members)
	}
}

func (h *ProjectHandler) NewMessage(c *gin.Context) {

	var dto dto.NewMessageDto

	if err := c.ShouldBindJSON(&dto); err != nil {

		c.Status(403)
		c.Abort()
		return
	}

	userID, _ := c.Get("userID")
	cardID := c.Param("card_id")

	dto.CardID = cardID
	dto.SenderID = userID.(string)

	if err := h.messageService.NewMessage(dto); err != nil {
		c.Status(403)
		c.Abort()
		return
	}

	c.Status(201)

}
