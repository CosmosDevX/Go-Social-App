// Package middleware
package middleware

import (
	"context"
	"myapp/internal/service/authorization"
	"net/http"
	"strings"
)

type UsernameContextKey struct{}
type UserIDContextKey struct{}

type AuthMiddleware struct {
	jwtService authorization.JWTServiceInterface
}

func NewAuthMiddleware(jwtService authorization.JWTServiceInterface) AuthMiddleware {
	return AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m AuthMiddleware) Protect(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "auth header is null", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "token is null", http.StatusUnauthorized)
			return
		}

		claims, err := m.jwtService.ParseToken(tokenString)
		if err != nil {
			http.Error(w, err.Message, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, UsernameContextKey{}, claims.Username)
		ctx = context.WithValue(ctx, UserIDContextKey{}, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
