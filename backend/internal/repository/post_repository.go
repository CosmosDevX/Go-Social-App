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

type PostRepositoryInterface interface {
	Create(postDTO dto.PostDTO, db *gorm.DB) (uint, *delivery.APIError)
	GetByID(postID uint, db *gorm.DB) (*model.Post, *delivery.APIError)
	GetAll(userID uint, db *gorm.DB) ([]model.Post, *delivery.APIError)
	IncrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError)
	DecrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError)
}

type PostRepository struct{}

func (r PostRepository) Create(postDTO dto.PostDTO, db *gorm.DB) (uint, *delivery.APIError) {
	post := model.Post{
		PostName:        postDTO.PostName,
		PostDescription: postDTO.PostDescription,
		CreatorID:       postDTO.CreatorID,
	}

	result := db.Create(&post)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.CreateError, Message: "error during post creating"}
	}

	return post.ID, nil
}

func (r PostRepository) GetByID(postID uint, db *gorm.DB) (*model.Post, *delivery.APIError) {
	var post model.Post

	result := db.First(&post, "id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "post not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding post by id"}
	}

	return &post, nil
}

func (r PostRepository) GetAll(userID uint, db *gorm.DB) ([]model.Post, *delivery.APIError) {
	var posts []model.Post

	result := db.Find(&posts, "creator_id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "posts not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding posts by user id"}
	}

	return posts, nil
}

func (r PostRepository) IncrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError) {
	var likes int

	result := db.Model(&model.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes + ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.UpdateError, Message: "error during increment likes on post"}
	}

	if result.RowsAffected == 0 {
		return 0, &delivery.APIError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}

func (r PostRepository) DecrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError) {
	var likes int
	result := db.Model(&model.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes - ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if result.Error != nil {
			if errors.Is(result.Error, context.DeadlineExceeded) {
				return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
			}

			return 0, &delivery.APIError{Code: constants.UpdateError, Message: "error during decrement likes on post"}
		}
	}

	if result.RowsAffected == 0 {
		return 0, &delivery.APIError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}
