package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenService interface {
	Validate(tokenString string) (uuid.UUID, error)
}

type tokenService struct {
	secret []byte
}

func NewTokenService(secret string) TokenService {
	return &tokenService{
		secret: []byte(secret),
	}
}

func (s *tokenService) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims")
	}
	userIDValue, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user id claim")
	}
	userID, err := uuid.Parse(userIDValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id: %w", err)
	}
	return userID, nil
}
