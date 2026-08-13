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
	"bugtracker/internal/security"
	"errors"
	"fmt"
)

type UserService struct {
	userRepo     *repositories.UserRepo
	tokenService *TokenService
}

func NewUserService(userRepo *repositories.UserRepo, tokenService *TokenService) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (u *UserService) Registration(dto *dto.UserRegistrationDto) (models.TokenReponse, error) {

	passwordHash, err := security.HashPassword(dto.Password)

	if err != nil {
		return models.TokenReponse{}, err
	}

	user := &models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Login:    dto.Login,
		Password: passwordHash,
	}

	u.userRepo.Create(user)

	return models.TokenReponse{}, nil

}

func (u *UserService) Login(dto *dto.UserLoginDto) (*models.TokenReponse, error) {

	fmt.Println("Start login...")

	user, err := u.userRepo.GetByLogin(dto.Login)

	fmt.Println("Got user")

	if err != nil {
		return &models.TokenReponse{}, err
	}

	verificationResult, _ := security.VerifyPassword(dto.Password, user.Password)

	fmt.Println("Verification result: ", verificationResult)

	if verificationResult == true {

		tokensPair, err := u.tokenService.GenerateTokensPair(user)

		if err != nil {
			return &models.TokenReponse{}, err
		}

		return tokensPair, nil

	}

	return &models.TokenReponse{}, errors.New("Unathorized")

}
