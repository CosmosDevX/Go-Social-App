// Package middleware
package middleware

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/logger"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
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
			utils.WriteError(w, *domain.NewDomainError(constants.AuthError, "auth header is null"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			utils.WriteError(w, *domain.NewDomainError(constants.InvalidTokenError, "token is null"))
			return
		}

		claims, err := m.jwtService.ParseToken(tokenString)
		if err != nil {
			logger.FromContext(r.Context()).Warn("auth middleware: invalid token", "code", err.Code)
			utils.WriteError(w, *domain.NewDomainError(constants.AuthError, err.Message))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UsernameContextKey{}, claims.Username)
		ctx = context.WithValue(ctx, UserIDContextKey{}, claims.UserID)
		ctx = logger.WithUserID(ctx, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
