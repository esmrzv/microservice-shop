package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

type mockTokenService struct {
	validateErr error
	userID      uuid.UUID
}

func (m mockTokenService) Validate(tokenString string) (uuid.UUID, error) {
	return m.userID, m.validateErr
}

func TestMiddleware_NoAuthorization(t *testing.T) {
	mock := &mockTokenService{}
	middleware := NewAuthMiddleware(mock)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireAuth(next)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	mock := &mockTokenService{
		validateErr: errors.New("invalid token"),
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	middleware := NewAuthMiddleware(mock)
	handler := middleware.RequireAuth(next)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	req.Header.Add("Authorization", "Bearer 123")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, rr.Code)

	}

}

func TestMiddleware_ValidToken(t *testing.T) {
	userID := uuid.New()
	mock := &mockTokenService{
		userID:      userID,
		validateErr: nil,
	}

	middleware := NewAuthMiddleware(mock)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUSerID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
		if !ok {
			t.Fatal("not found user in context")
		}
		if gotUSerID != userID {
			t.Errorf("expected userID %s, got %s", userID, gotUSerID)
		}
	})
	handler := middleware.RequireAuth(next)
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	req.Header.Add("Authorization", "Bearer 123")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

}
