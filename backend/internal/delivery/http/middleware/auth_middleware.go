// Package middleware
package middleware

import (
	"context"
	"myapp/internal/service/authorization"
	"net/http"
	"strings"
)

type UserContextKey struct{}

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

		ctx := context.WithValue(r.Context(), UserContextKey{}, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
