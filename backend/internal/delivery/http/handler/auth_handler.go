// Package handler
package handler

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
	"net/http"
)

type AuthService interface {
	Auth(ctx context.Context, userDTO dto.UserDTO) (*authorization.AuthResult, *domain.DomainError)
	Refresh(ctx context.Context, oldRefreshToken string) (*authorization.AuthResult, *domain.DomainError)
	Logout(ctx context.Context, userID int, refreshToken string) *domain.DomainError
}

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

func (h AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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
