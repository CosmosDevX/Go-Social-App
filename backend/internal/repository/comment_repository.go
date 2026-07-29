package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type CommentRepositoryInterface interface {
	Delete(commentID, userID uint, db *gorm.DB) *domain.DomainError
	Create(commentDTO dto.CommentDTO, db *gorm.DB) (uint, *domain.DomainError)
	GetAllByPostID(postID uint, db *gorm.DB) ([]domain.Comment, *domain.DomainError)
	CountCommentsOnPost(postID uint, db *gorm.DB) (int, *domain.DomainError)
}

type CommentRepository struct{}

func (r CommentRepository) Create(commentDTO dto.CommentDTO, db *gorm.DB) (uint, *domain.DomainError) {
	comment := domain.Comment{
		CommentText: commentDTO.CommentText,
		PostID:      commentDTO.PostID,
		CreatorID:   commentDTO.CreatorID,
	}

	result := db.Create(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during comment creating"}
	}

	return comment.ID, nil
}

func (r CommentRepository) GetAllByPostID(postID uint, db *gorm.DB) ([]domain.Comment, *domain.DomainError) {
	var comments []domain.Comment

	result := db.Joins("Creator").Find(&comments, "post_id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "comments not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during finding comments by post id"}
	}

	return comments, nil
}

func (r CommentRepository) CountCommentsOnPost(postID uint, db *gorm.DB) (int, *domain.DomainError) {
	var count int64
	result := db.Model(&domain.Comment{}).Where("post_id = ?", postID).Count(&count)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.FindError, Message: "error during count comments on post"}
	}

	return int(count), nil
}

func (r CommentRepository) Delete(commentID, userID uint, db *gorm.DB) *domain.DomainError {
	result := db.Delete(&domain.Comment{}, "id = ? AND creator_id = ?", commentID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &domain.DomainError{Code: constants.DeleteError, Message: "error during deleting the comment"}
	}

	if result.RowsAffected == 0 {
		return &domain.DomainError{Code: constants.NotFound, Message: "comment not deleted"}
	}

	return nil
}
