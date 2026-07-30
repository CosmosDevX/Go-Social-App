// Package handler
package handler

import (
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type AuthHandler struct {
	authService authorization.AuthServiceInterface
}

func NewAuthHandler(authService authorization.AuthServiceInterface) AuthHandler {
	return AuthHandler{
		authService: authService,
	}
}

func (h AuthHandler) HandleAuth(w http.ResponseWriter, r *http.Request) {
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

func (h AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
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

func (h AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenCookie, err := r.Cookie(constants.RefreshTokenKey)
	if err != nil || tokenCookie.Value == "" {
		utils.WriteJSON(w, map[string]string{"message": "refresh token not exists"})
		return
	}

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	if domainErr := h.authService.Logout(strconv.Itoa(int(userID)), tokenCookie.Value, ctx); domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "logout successful"})
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
