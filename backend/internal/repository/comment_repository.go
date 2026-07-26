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

type CommentRepositoryInterface interface {
	Create(commentDTO dto.CommentDTO, db *gorm.DB) (uint, *delivery.APIError)
	GetAllByPostID(postID uint, db *gorm.DB) ([]model.Comment, *delivery.APIError)
}

type CommentRepository struct{}

func (r CommentRepository) Create(commentDTO dto.CommentDTO, db *gorm.DB) (uint, *delivery.APIError) {
	comment := model.Comment{
		CommentText: commentDTO.CommentText,
		PostID:      commentDTO.PostID,
		CreatorID:   commentDTO.CreatorID,
	}

	result := db.Create(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.CreateError, Message: "error during comment creating"}
	}

	return comment.ID, nil
}

func (r CommentRepository) GetAllByPostID(postID uint, db *gorm.DB) ([]model.Comment, *delivery.APIError) {
	var comments []model.Comment

	result := db.Joins("Creator").Find(&comments, "post_id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "comments not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding comments by post id"}
	}

	return comments, nil
}
