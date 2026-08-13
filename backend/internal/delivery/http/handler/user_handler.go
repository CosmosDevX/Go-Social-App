package handler

import (
	"context"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

type UserService interface {
	CreateUser(ctx context.Context, userDTO dto.UserDTO) (int, *domain.DomainError)
	CurrentUserProfile(ctx context.Context, userID int) (string, *domain.DomainError)
	GetUsernameByID(ctx context.Context, mpathValue string) (string, *domain.DomainError)
}

type UserHandler struct {
	userService UserService
	rateLimiter redis_rate.Limiter
}

func NewUserHandler(userService UserService, rateLimiter redis_rate.Limiter) UserHandler {
	return UserHandler{
		userService: userService,
		rateLimiter: rateLimiter,
	}
}

// HandleCreateUser godoc
// @Summary      Регистрация нового пользователя
// @Description  Создаёт пользователя. Username 3-60 символов, password 10-100.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserDTO  true  "Данные пользователя"
// @Success      200  {object}  UserIDResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      409  {object}  ErrorResponse  "UNIQUE_VIOLATION"
// @Router       /user/create [post]
func (h UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	if err := utils.ActivateRateLimiter(ctx, w, r, "create_user", &h.rateLimiter, redis_rate.PerHour(4)); err != nil {
		return
	}

	id, domainErr := h.userService.CreateUser(ctx, userDTO)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"user_id": id})
}

// HandleCurrentUserProfile godoc
// @Summary      Профиль текущего пользователя
// @Description  Возвращает username авторизованного пользователя
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  UsernameResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /user/current/profile [get]
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

// HandleGetUsernameByID godoc
// @Summary      Получить username по ID
// @Description  Возвращает username пользователя по его ID
// @Tags         users
// @Produce      json
// @Param        user_id  path  int  true  "ID пользователя"
// @Success      200  {object}  UsernameResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /user/get_username_by_id/{user_id} [get]
func (h UserHandler) HandleGetUsernameByID(w http.ResponseWriter, r *http.Request) {
	username, domainErr := h.userService.GetUsernameByID(r.Context(), r.PathValue("user_id"))
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"username": username})
}
