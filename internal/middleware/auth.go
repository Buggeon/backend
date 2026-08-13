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

package middleware

import (
	"bugtracker/internal/services"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService *services.TokenService
}

func NewAuthMiddleware(tokenService *services.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {

	fmt.Println("Welcome to the middleware")

	return func(ctx *gin.Context) {

		fmt.Println("You're in middleware")

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.Status(401)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Status(401)
			ctx.Abort()
			return
		}

		token := parts[1]
		claims, err := m.tokenService.ValidateAccessToken(token)

		if err != nil {
			ctx.Status(401)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Next()
	}
}
