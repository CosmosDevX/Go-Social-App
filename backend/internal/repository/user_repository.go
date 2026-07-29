// Package repository
package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByName(username string, db *gorm.DB) (*domain.User, *domain.DomainError)
	CreateUser(userDTO dto.UserDTO, db *gorm.DB) (uint, *domain.DomainError)
	GetUsernameByID(userID uint, db *gorm.DB) (string, *domain.DomainError)
}

type UserRepository struct{}

func (r UserRepository) GetUserByName(username string, db *gorm.DB) (*domain.User, *domain.DomainError) {
	var user domain.User
	result := db.First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "user not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during finding user by name"}
	}

	return &user, nil
}

func (r UserRepository) CreateUser(userDTO dto.UserDTO, db *gorm.DB) (uint, *domain.DomainError) {
	user := domain.User{Username: userDTO.Username, Password: userDTO.Password}
	result := db.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, &domain.DomainError{Code: constants.UniqueViolation, Message: "username already taken"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during user creating"}
	}

	return user.ID, nil
}

func (r UserRepository) GetUsernameByID(userID uint, db *gorm.DB) (string, *domain.DomainError) {
	var username string
	result := db.Model(&domain.User{}).Where("id = ?", userID).Select("username").Scan(&username)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return "", &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", &domain.DomainError{Code: constants.NotFound, Message: "username not found not found"}
		}

		return "", &domain.DomainError{Code: constants.FindError, Message: "error during find username by user id"}
	}

	return username, nil
}
