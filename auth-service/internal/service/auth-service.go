package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/esmrzv/microservice-shop/auth-service/internal/auth"
	"github.com/esmrzv/microservice-shop/auth-service/internal/models"
	"github.com/esmrzv/microservice-shop/auth-service/internal/repository"
	"github.com/esmrzv/microservice-shop/auth-service/internal/service/dto"
	"golang.org/x/crypto/bcrypt"
)


var ErrInvalidCredentials = errors.New("ivalid credentials")

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	
}


type authService struct {
	repo repository.UserRepository
	token auth.TokenService
}

func NewAuthService(r repository.UserRepository, t auth.TokenService) AuthService {
	return &authService{
		repo: r,
		token: t,
	}
	
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error){
	
	_, err := s.repo.GetByEmail(ctx, req.Email)
	if err == nil{
		return nil, fmt.Errorf("userwith email: %s already exists", req.Email)
	}
	if !errors.Is(err, repository.ErrUserNotFound){
		 return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil{
		return nil, err
	}

	user := &models.User{
		Email: req.Email,
		PasswordHash: string(hashedPassword),
	}
	err = s.repo.CreateUser(ctx, user)
	if err != nil{
		return nil, fmt.Errorf("failed to create user: %w", err)
	} 

	return &dto.RegisterResponse{
		Email: user.Email,
	}, nil
    
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error){
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrUserNotFound){
		return nil, ErrInvalidCredentials
	}
	if err != nil{
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil{
		return nil, errors.New("invalid credentials")
	}
	token, err := s.token.Generate(user.ID)
	if err != nil{
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return &dto.LoginResponse{
		Token: token,
	}, nil

}