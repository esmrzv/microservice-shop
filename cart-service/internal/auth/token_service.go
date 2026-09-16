package auth

import (
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
	if !token.Valid {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, err
	}

	userIdValue, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(userIdValue)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
