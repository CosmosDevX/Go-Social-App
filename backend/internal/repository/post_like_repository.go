package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type PostLikeRepository struct{}

func (r PostLikeRepository) GetLikedPostsID(userID uint, db *gorm.DB) ([]uint, *domain.DomainError) {
	var postIDs []uint

	result := db.Model(&domain.PostLike{}).
		Where("liked_user_id = ?", userID).
		Pluck("post_id", &postIDs)

	if result.Error != nil {
		return nil, &domain.DomainError{
			Code:    constants.FindError,
			Message: "error during finding liked post ids",
		}
	}

	return postIDs, nil
}

func (r PostLikeRepository) LikeExists(userID, postID uint, db *gorm.DB) (bool, *domain.DomainError) {
	var postLike domain.PostLike
	result := db.Model(&domain.PostLike{}).Where("liked_user_id = ? AND post_id = ?", userID, postID).First(&postLike)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return false, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, &domain.DomainError{Code: constants.FindError, Message: "error during finding like on post"}
	}

	return true, nil
}

func (r PostLikeRepository) CreateLike(likedUserID uint, postID uint, db *gorm.DB) *domain.DomainError {
	postLike := domain.PostLike{
		LikedUserID: likedUserID,
		PostID:      postID,
	}

	result := db.Create(&postLike)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &domain.DomainError{Code: constants.CreateError, Message: "error during post creating"}
	}

	return nil
}

func (r PostLikeRepository) DeleteLike(likedUserID, postID uint, db *gorm.DB) *domain.DomainError {
	result := db.Delete(&domain.PostLike{}, "liked_user_id = ? AND post_id = ?", likedUserID, postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &domain.DomainError{Code: constants.DeleteError, Message: "error during post like deleting"}
	}

	if result.RowsAffected == 0 {
		return &domain.DomainError{Code: constants.NotFound, Message: "this user like on post not found"}
	}

	return nil
}
