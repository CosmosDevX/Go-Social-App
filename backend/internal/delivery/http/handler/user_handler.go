package handler

import (
	"context"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserDTO) (int, *domain.DomainError)
	CurrentUserProfile(ctx context.Context, userID int) (string, *domain.DomainError)
	GetUsernameByID(ctx context.Context, mpathValue string) (string, *domain.DomainError)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

func (h UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	id, domainErr := h.userService.CreateUser(r.Context(), userDTO)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"user_id": id})
}

func (h UserHandler) HandleCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	username, domainErr := h.userService.CurrentUserProfile(ctx, parsedUserID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"username": username})
}

func (h UserHandler) HandleGetUsernameByID(w http.ResponseWriter, r *http.Request) {
	username, domainErr := h.userService.GetUsernameByID(r.Context(), r.PathValue("user_id"))
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"username": username})
}
