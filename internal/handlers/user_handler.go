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

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Registration(c *gin.Context) {

	var dto dto.UserRegistrationDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(403)
	}

	h.userService.Registration(&dto)

}

func (h *UserHandler) Login(c *gin.Context) {

	var dto dto.UserLoginDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Status(403)
	}

	response, err := h.userService.Login(&dto)

	if err != nil {
		c.Status(401)
	} else {
		c.JSON(200, response)
	}

}
