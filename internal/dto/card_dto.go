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

package dto

type CreateCardDto struct {
	Title     string   `json:"title" binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Assignees []string `json:"assignees" binding:"required"`
	BoardID   string   `json:"boardId"`
	Priority  string   `json:"priority" binding:"required"`
}

type GetCard struct {
	CardID string `json:"cardId"  binding:"required"`
}

type GetCards struct {
	BoardID string `json:"boardId" binding:"required"`
}

type DeleteCard struct {
	CardID string `json:"cardId" binding:"required"`
}
