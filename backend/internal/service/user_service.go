package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/utils"
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

func (s UserService) CreateUser(userDTO dto.UserDTO, ctx context.Context) (int, *domain.DomainError) {
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

func (s UserService) CurrentUserProfile(ctx context.Context) (string, *domain.DomainError) {
	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		return "", &domain.DomainError{Code: constants.ParseError, Message: parseErr.Error()}
	}

	username, domainErr := s.userRepository.GetUsernameByID(ctx, int(parsedUserID))
	if domainErr != nil {
		return "", domainErr
	}

	return username, nil
}

func (s UserService) GetUsernameByID(pathValue string, ctx context.Context) (string, *domain.DomainError) {
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
