// Package repository
package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByName(username string, db *gorm.DB) (*model.User, *delivery.APIError)
	CreateUser(userDTO dto.UserDTO, db *gorm.DB) (uint, *delivery.APIError)
	GetUsernameByID(userID uint, db *gorm.DB) (string, *delivery.APIError)
}

type UserRepository struct{}

func (r UserRepository) GetUserByName(username string, db *gorm.DB) (*model.User, *delivery.APIError) {
	var user model.User
	result := db.First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "user not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding user by name"}
	}

	return &user, nil
}

func (r UserRepository) CreateUser(userDTO dto.UserDTO, db *gorm.DB) (uint, *delivery.APIError) {
	user := model.User{Username: userDTO.Username, Password: userDTO.Password}
	result := db.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.CreateError, Message: "error during user creating"}
	}

	return user.ID, nil
}

func (r UserRepository) GetUsernameByID(userID uint, db *gorm.DB) (string, *delivery.APIError) {
	var username string
	result := db.Model(&model.User{}).Where("id = ?", userID).Select("username").Scan(&username)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return "", &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", &delivery.APIError{Code: constants.NotFound, Message: "username not found not found"}
		}

		return "", &delivery.APIError{Code: constants.FindError, Message: "error during find username by user id"}
	}

	return username, nil
}
