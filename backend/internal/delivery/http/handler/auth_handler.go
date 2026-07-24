// Package handler
package handler

import (
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
	"net/http"
)

type AuthHandlerInterface interface {
	AuthHandler(w http.ResponseWriter, r *http.Request)
	RefreshHandler(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	authService authorization.AuthServiceInterface
}

func NewAuthHandler(authService authorization.AuthServiceInterface) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

func (h AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userDTO dto.UserDTO
	if err := utils.Deserialize(r.Body, &userDTO); err != nil {
		http.Error(w, "error during deserialization user", http.StatusBadRequest)
		return
	}

	validationErr := userDTO.Validate()
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return
	}

	authResult, err := h.authService.Auth(userDTO, ctx)
	if err != nil {
		http.Error(w, err.Message, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		http.Error(w, "refresh token not exists in cookie", http.StatusUnauthorized)
		return
	}

	authResult, apiErr := h.authService.Refresh(tokenCookie.Value, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
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
