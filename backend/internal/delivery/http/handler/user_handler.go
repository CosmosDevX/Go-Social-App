package handler

import (
	"myapp/internal/delivery/http/dto"
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

func (h UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO dto.UserDTO
	if err := utils.Deserialize(r.Body, &userDTO); err != nil {
		http.Error(w, "error during deserializing body", http.StatusBadRequest)
		return
	}

	validationErr := userDTO.Validate()
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.userService.CreateUser(&userDTO, r.Context())
	if err != nil {
		http.Error(w, err.Message, utils.IdentifyRepositoryError(err.Code))
		return
	}

	utils.WriteJSON(w, map[string]uint{"user_id": id})
}
