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
	Delete(commentID, userID uint, db *gorm.DB) *delivery.APIError
	Create(commentDTO dto.CommentDTO, db *gorm.DB) (uint, *delivery.APIError)
	GetAllByPostID(postID uint, db *gorm.DB) ([]model.Comment, *delivery.APIError)
	CountCommentsOnPost(postID uint, db *gorm.DB) (int, *delivery.APIError)
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

func (r CommentRepository) CountCommentsOnPost(postID uint, db *gorm.DB) (int, *delivery.APIError) {
	var count int64
	result := db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.FindError, Message: "error during count comments on post"}
	}

	return int(count), nil
}

func (r CommentRepository) Delete(commentID, userID uint, db *gorm.DB) *delivery.APIError {
	result := db.Delete(&model.Comment{}, "id = ? AND creator_id = ?", commentID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &delivery.APIError{Code: constants.DeleteError, Message: "error during deleting the comment"}
	}

	if result.RowsAffected == 0 {
		return &delivery.APIError{Code: constants.NotFound, Message: "comment not deleted"}
	}

	return nil
}
