package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func TestTokenService_Validate(t *testing.T) {
	secret := "test-secret"
	service := NewTokenService(secret)

	userID := uuid.New()

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	gotUserID, err := service.Validate(tokenString)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gotUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, gotUserID)
	}
}

func TestTokenService_Validate_WrongSecret(t *testing.T) {
	serviceSecret := "wrong-secret"
	tokenSecret := "correct-secret"
	service := NewTokenService(serviceSecret)

	userID := uuid.New()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = service.Validate(tokenString)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

}

func TestUUIDParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "valid Uuid",
			input: uuid.New().String(),
			valid: true,
		},
		{
			name:  "invalid uuid",
			input: "",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uuid.Parse(tt.input)
			if tt.valid && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected error, got nil")
			}
		})
	}
}

func TestTokenService_Validate_TableDriven(t *testing.T) {
	serviceSecret := "test-secret"
	service := NewTokenService(serviceSecret)

	tests := []struct {
		name    string
		userID  string
		exp     time.Time
		wantErr bool
	}{
		{
			name:    "valid token",
			userID:  uuid.New().String(),
			exp:     time.Now().Add(time.Hour),
			wantErr: false,
		},
		{
			name:    "missing userId",
			userID:  "",
			exp:     time.Now().Add(time.Hour),
			wantErr: true,
		},
		{
			name:    "invalid userId",
			userID:  "hello",
			exp:     time.Now().Add(time.Hour),
			wantErr: true,
		},
		{
			name:    "expired token",
			userID:  uuid.New().String(),
			exp:     time.Now().Add(-time.Hour),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"user_id": tt.userID,
				"iat":     time.Now().Unix(),
				"exp":     tt.exp.Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(serviceSecret))
			if err != nil {
				t.Fatalf("failed to sign token: %v", err)
			}
			_, err = service.Validate(tokenString)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
