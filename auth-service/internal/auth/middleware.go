package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	userIdKey contextKey = "userID"
)

type AuthMiddleware interface {
	RequireAuth(next http.Handler) http.Handler
}

type authMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret string) AuthMiddleware {
	return &authMiddleware{
		secret: []byte(secret),
	}
}

func (m *authMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Fields(authorization)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return m.secret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "invalid user_id", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), userIdKey, userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GetUserId(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(userIdKey).(string)
	return userId, ok
}
