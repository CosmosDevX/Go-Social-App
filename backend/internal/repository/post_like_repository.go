package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type PostLikeRepositoryInterface interface {
	CreateLike(likedUserID uint, postID uint, db *gorm.DB) *delivery.APIError
	DeleteLike(likedUserID, postID uint, db *gorm.DB) *delivery.APIError
	LikeExists(userID, postID uint, db *gorm.DB) (bool, *delivery.APIError)
	GetLikedPostsID(userID uint, db *gorm.DB) ([]uint, *delivery.APIError)
}

type PostLikeRepository struct{}

func (r PostLikeRepository) GetLikedPostsID(userID uint, db *gorm.DB) ([]uint, *delivery.APIError) {
	var postIDs []uint

	result := db.Model(&model.PostLike{}).
		Where("liked_user_id = ?", userID).
		Pluck("post_id", &postIDs)

	if result.Error != nil {
		return nil, &delivery.APIError{
			Code:    constants.FindError,
			Message: "error during finding liked post ids",
		}
	}

	return postIDs, nil
}

func (r PostLikeRepository) LikeExists(userID, postID uint, db *gorm.DB) (bool, *delivery.APIError) {
	var postLike model.PostLike
	result := db.Model(&model.PostLike{}).Where("liked_user_id = ? AND post_id = ?", userID, postID).First(&postLike)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return false, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, &delivery.APIError{Code: constants.FindError, Message: "error during finding like on post"}
	}

	return true, nil
}

func (r PostLikeRepository) CreateLike(likedUserID uint, postID uint, db *gorm.DB) *delivery.APIError {
	postLike := model.PostLike{
		LikedUserID: likedUserID,
		PostID:      postID,
	}

	result := db.Create(&postLike)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &delivery.APIError{Code: constants.CreateError, Message: "error during post creating"}
	}

	return nil
}

func (r PostLikeRepository) DeleteLike(likedUserID, postID uint, db *gorm.DB) *delivery.APIError {
	result := db.Delete(&model.PostLike{}, "liked_user_id = ? AND post_id = ?", likedUserID, postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &delivery.APIError{Code: constants.DeleteError, Message: "error during post like deleting"}
	}

	if result.RowsAffected == 0 {
		return &delivery.APIError{Code: constants.NotFound, Message: "this user like on post not found"}
	}

	return nil
}
