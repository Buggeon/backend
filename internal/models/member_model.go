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

type Member struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	UserId     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Role       string             `bson:"role" json:"role"`
	ProjectID  primitive.ObjectID `bson:"project_id" json:"project_id"`
	Directions []string           `bson:"directions" json:"directions"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
