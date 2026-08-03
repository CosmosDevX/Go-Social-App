// Package handler
package handler

import (
	"context"
	"fmt"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type AuthService interface {
	Auth(ctx context.Context, userDTO dto.UserDTO) (*authorization.AuthResult, *domain.DomainError)
	Refresh(ctx context.Context, oldRefreshToken string) (*authorization.AuthResult, *domain.DomainError)
	Logout(ctx context.Context, userID int, refreshToken string) *domain.DomainError
}

type AuthHandler struct {
	authService AuthService
	rateLimiter redis_rate.Limiter
}

func NewAuthHandler(authService AuthService, rateLimiter redis_rate.Limiter) AuthHandler {
	return AuthHandler{
		authService: authService,
		rateLimiter: rateLimiter,
	}
}

func (h AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.RemoteAddr

	res, err := h.rateLimiter.Allow(ctx, "auth:"+ip, redis_rate.PerMinute(5))
	if err != nil {
		WriteError(w, *domain.NewDomainError(constants.TooManyRequests, "too many requests"))
		return
	}

	if res.Allowed == 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", res.RetryAfter.Seconds()))
		WriteError(w, *domain.NewDomainError(constants.TooManyRequests, "too many requests. Try again later"))
		return
	}

	var userDTO dto.UserDTO
	if err := utils.Deserialize(r.Body, &userDTO); err != nil {
		WriteError(w, *domain.NewDeserializingError("error during deserializing user"))
		return
	}

	validationErr := userDTO.Validate()
	if validationErr != nil {
		WriteError(w, *domain.NewValidationError(validationErr.Error()))
		return
	}

	authResult, domainErr := h.authService.Auth(ctx, userDTO)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		WriteError(w, *domain.NewTokenError("refresh token not exists"))
		return
	}

	authResult, domainErr := h.authService.Refresh(ctx, tokenCookie.Value)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		WriteJSON(w, map[string]string{"message": "refresh token not exists"})
		return
	}

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	if domainErr := h.authService.Logout(ctx, userID, tokenCookie.Value); domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, map[string]string{"message": "logout successful"})
}

func (h AuthHandler) newRefreshTokenCookie(refreshToken string) *http.Cookie {
	return &http.Cookie{
		Name:     constants.RefreshTokenKey,
		Value:    refreshToken,
		MaxAge:   constants.RefreshTokenMaxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}
