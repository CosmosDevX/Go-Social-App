package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	CreateUser(userDTO *dto.UserDTO, ctx context.Context) (uint, *delivery.APIError)
}

type UserService struct {
	userRepository repository.UserRepositoryInterface
}

func NewUserService(userRepository repository.UserRepositoryInterface) UserServiceInterface {
	return UserService{
		userRepository: userRepository,
	}
}

func (s UserService) CreateUser(userDTO *dto.UserDTO, ctx context.Context) (uint, *delivery.APIError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 10)
	if err != nil {
		return 0, &delivery.APIError{Code: constants.InvalidPassword, Message: "error during password hashing"}
	}
	userDTO.Password = string(hashedPassword)

	id, apiErr := s.userRepository.CreateUser(*userDTO, ctx)
	if apiErr != nil {
		return 0, apiErr
	}

	return id, nil
}
