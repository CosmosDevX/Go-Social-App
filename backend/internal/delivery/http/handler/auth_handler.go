// Package handler
package handler

import (
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
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
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing user"))
		return
	}

	validationErr := userDTO.Validate()
	if validationErr != nil {
		utils.WriteError(w, *domain.NewValidationError(validationErr.Error()))
		return
	}

	authResult, domainErr := h.authService.Auth(userDTO, ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	http.SetCookie(w, h.newRefreshTokenCookie(authResult.RefreshToken))

	utils.WriteJSON(w, map[string]string{"access_token": authResult.AccessToken})
}

func (h AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteError(w, *domain.NewTokenError("refresh token not exists"))
		return
	}

	authResult, domainErr := h.authService.Refresh(tokenCookie.Value, ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
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
