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
	"bugtracker/config"
	"bugtracker/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	cfg config.Config
}

func NewTokenService(cfg config.Config) *TokenService {
	return &TokenService{
		cfg: cfg,
	}
}

type TokenClaims struct {
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	UserLogin string `json:"user_login"`
	UserEmail string `json:"user_email"`
	jwt.RegisteredClaims
}

func (s *TokenService) GenerateAccessToken(user *models.User) (string, error) {

	expiry := time.Now().Add(time.Duration(s.cfg.AccessTokenExpire) * time.Minute)

	claims := TokenClaims{
		UserID:    user.ID.Hex(),
		UserName:  user.Name,
		UserLogin: user.Login,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.TokensIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.AccessTokenSecret))

}

func (s *TokenService) GenerateRefreshToken(user *models.User) (string, error) {

	expiry := time.Now().Add(time.Duration(s.cfg.RefreshTokenExpire) * time.Hour)

	claims := TokenClaims{
		UserID:    user.ID.Hex(),
		UserName:  user.Name,
		UserLogin: user.Login,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.TokensIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.RefreshTokenSecret))

}

func (s *TokenService) ValidateAccessToken(accessToken string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(s.cfg.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")

}

func (s *TokenService) ValidateRefreshToken(refreshToken string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(s.cfg.RefreshTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")

}

func (s *TokenService) RefreshAccessToken(refreshToken string, user *models.User) (string, error) {

	_, err := s.ValidateRefreshToken(refreshToken)

	if err != nil {
		return "", errors.New("Invalid token")
	}

	return s.GenerateAccessToken(user)

}

func (s *TokenService) GenerateTokensPair(user *models.User) (models.TokenReponse, error) {

	accessToken, err := s.GenerateAccessToken(user)

	fmt.Println(err)

	if err != nil {
		return models.TokenReponse{}, nil
	}

	refreshToken, err := s.GenerateRefreshToken(user)

	if err != nil {
		return models.TokenReponse{}, nil
	}

	return models.TokenReponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
	}, nil
}
