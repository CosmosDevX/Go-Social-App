package handler

import (
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/internal/utils"
	"net/http"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userService service.UserServiceInterface) UserHandler {
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

	id, domainErr := h.userService.CreateUser(&userDTO, r.Context())
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]uint{"user_id": id})
}

func (h UserHandler) HandleCurrentUserProfile(w http.ResponseWriter, r *http.Request) {
	username, domainErr := h.userService.CurrentUserProfile(r.Context())
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"username": username})
}

func (h UserHandler) HandleGetUsernameByID(w http.ResponseWriter, r *http.Request) {
	username, domainErr := h.userService.GetUsernameByID(r.PathValue("user_id"), r.Context())
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"username": username})
}
