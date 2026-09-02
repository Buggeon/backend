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

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title     string               `bson:"title" json:"title"`
	Content   string               `bson:"content" json:"content"`
	Assignees []primitive.ObjectID `bson:"assignees" json:"assignees"`
	BoardID   primitive.ObjectID   `bson:"board_id" json:"board_id"`
	Priority  string               `bson:"priority" json:"priority"`
	Status    string               `bson:"status" json:"status"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	Messages  []primitive.ObjectID `bson:"messages" json:"messages"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}
