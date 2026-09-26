package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/krishnasingh0210/blog/pkg/jwt"
	"github.com/krishnasingh0210/blog/pkg/response"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
	ContextRoleKey   contextKey = "role"
)

func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing or invalid authorization header")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.Verify(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ContextRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}