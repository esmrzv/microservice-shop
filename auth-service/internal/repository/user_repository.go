package repository

import (
	"context"
	"errors"

	"github.com/esmrzv/microservice-shop/auth-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)
	


type UserRepository interface{
	CreateUser(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

var ErrUserNotFound = errors.New("user not found")

type userRepository struct {
	db *pgxpool.Pool
	
}


func NewUserRepository(db *pgxpool.Pool) UserRepository{
	return &userRepository{db: db}

}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users(email, password_hash)
		VALUES($1, $2)
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query, 
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID)
	return err

}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error){
	query := `
		SELECT id, email, password_hash FROM users
		WHERE email= $1;
	`
	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil{
		if errors.Is(err, pgx.ErrNoRows){
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
	
}
