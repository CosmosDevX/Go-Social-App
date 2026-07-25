package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/repository"
	"myapp/internal/utils"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	CreateUser(userDTO *dto.UserDTO, ctx context.Context) (uint, *delivery.APIError)
	CurrentUserProfile(ctx context.Context) (string, *delivery.APIError)
	GetUsernameByID(pathValue string, ctx context.Context) (string, *delivery.APIError)
}

type UserService struct {
	userRepository repository.UserRepositoryInterface
	db             *gorm.DB
}

func NewUserService(userRepository repository.UserRepositoryInterface, db *gorm.DB) UserServiceInterface {
	return UserService{
		userRepository: userRepository,
		db:             db,
	}
}

func (s UserService) CreateUser(userDTO *dto.UserDTO, ctx context.Context) (uint, *delivery.APIError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 10)
	if err != nil {
		return 0, &delivery.APIError{Code: constants.InvalidPassword, Message: "error during password hashing"}
	}
	userDTO.Password = string(hashedPassword)

	id, apiErr := s.userRepository.CreateUser(*userDTO, s.db.WithContext(ctx))
	if apiErr != nil {
		return 0, apiErr
	}

	return id, nil
}

func (s UserService) CurrentUserProfile(ctx context.Context) (string, *delivery.APIError) {
	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		return "", &delivery.APIError{Code: constants.ParseError, Message: parseErr.Error()}
	}

	username, apiErr := s.userRepository.GetUsernameByID(parsedUserID, s.db.WithContext(ctx))
	if apiErr != nil {
		return "", apiErr
	}

	return username, nil
}

func (s UserService) GetUsernameByID(pathValue string, ctx context.Context) (string, *delivery.APIError) {
	userID, parseErr := strconv.ParseUint(pathValue, 10, 64)
	if parseErr != nil {
		return "", &delivery.APIError{Code: constants.ParseError, Message: "error during parsing user id"}
	}

	username, apiErr := s.userRepository.GetUsernameByID(uint(userID), s.db.WithContext(ctx))
	if apiErr != nil {
		return "", apiErr
	}

	return username, nil
}
