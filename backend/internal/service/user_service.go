package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError)
	CreateUser(ctx context.Context, userDTO dto.UserDTO) (int, *domain.DomainError)
	GetUsernameByID(ctx context.Context, userID int) (string, *domain.DomainError)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (s UserService) CreateUser(ctx context.Context, userDTO dto.UserDTO) (int, *domain.DomainError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 10)
	if err != nil {
		return 0, &domain.DomainError{Code: constants.InvalidPassword, Message: "error during password hashing"}
	}
	userDTO.Password = string(hashedPassword)

	id, domainErr := s.userRepository.CreateUser(ctx, userDTO)
	if domainErr != nil {
		return 0, domainErr
	}

	return id, nil
}

func (s UserService) CurrentUserProfile(ctx context.Context, userID int) (string, *domain.DomainError) {
	username, domainErr := s.userRepository.GetUsernameByID(ctx, userID)
	if domainErr != nil {
		return "", domainErr
	}

	return username, nil
}

func (s UserService) GetUsernameByID(ctx context.Context, pathValue string) (string, *domain.DomainError) {
	userID, parseErr := strconv.ParseInt(pathValue, 10, 64)
	if parseErr != nil {
		return "", &domain.DomainError{Code: constants.ParseError, Message: "error during parsing user id"}
	}

	username, domainErr := s.userRepository.GetUsernameByID(ctx, int(userID))
	if domainErr != nil {
		return "", domainErr
	}

	return username, nil
}
