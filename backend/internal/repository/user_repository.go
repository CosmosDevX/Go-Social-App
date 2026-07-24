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
	GetUserByName(username string, ctx context.Context) (*model.User, *delivery.APIError)
	CreateUser(userDTO dto.UserDTO, ctx context.Context) (uint, *delivery.APIError)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return UserRepository{
		db: db,
	}
}

func (r UserRepository) GetUserByName(username string, ctx context.Context) (*model.User, *delivery.APIError) {
	var user model.User
	result := r.db.WithContext(ctx).First(&user, "username = ?", username)
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

func (r UserRepository) CreateUser(userDTO dto.UserDTO, ctx context.Context) (uint, *delivery.APIError) {
	user := model.User{Username: userDTO.Username, Password: userDTO.Password}
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.CreateError, Message: "error during user creating"}
	}

	return user.ID, nil
}
