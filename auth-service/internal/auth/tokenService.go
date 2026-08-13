package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)


type TokenService interface {
	Generate(userID uuid.UUID) (string, error)
}


type tokenService struct {
	secret []byte
}

func NewTokenService(secret string) TokenService {
	return &tokenService{
		secret: []byte(secret),
	}
}


func (s *tokenService) Generate(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(s.secret)
	if err != nil{
		return "", err
	}
	return signedToken, nil

}