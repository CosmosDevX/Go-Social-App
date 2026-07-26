package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/repository"

	"gorm.io/gorm"
)

type PostLikeServiceInterface interface {
	LikePost(postID, userID uint, ctx context.Context) (int, *delivery.APIError)
	DislikePost(postID, userID uint, ctx context.Context) (int, *delivery.APIError)
}

type PostLikeService struct {
	postRepository     repository.PostRepositoryInterface
	postLikeRepository repository.PostLikeRepositoryInterface
	db                 *gorm.DB
}

func NewPostLikeService(postRepository repository.PostRepositoryInterface, postLikeRepository repository.PostLikeRepositoryInterface, db *gorm.DB) PostLikeServiceInterface {
	return PostLikeService{
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		db:                 db,
	}
}

func (s PostLikeService) LikePost(postID, userID uint, ctx context.Context) (int, *delivery.APIError) {
	dbWithCtx := s.db.WithContext(ctx)

	result, apiErr := s.postLikeRepository.LikeExists(userID, postID, dbWithCtx)
	if apiErr != nil {
		return 0, apiErr
	}
	if result {
		return 0, &delivery.APIError{Code: constants.CreateError, Message: "user already liked this post"}
	}

	tx := dbWithCtx.Begin()
	if tx.Error != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "error starting transaction"}
	}

	likes, apiErr := s.postRepository.IncrementLikes(postID, tx)
	if apiErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, apiErr
	}

	if apiErr := s.postLikeRepository.CreateLike(userID, postID, tx); apiErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, apiErr
	}

	if tx.Commit().Error != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return likes, nil
}

func (s PostLikeService) DislikePost(postID, userID uint, ctx context.Context) (int, *delivery.APIError) {
	dbWithCtx := s.db.WithContext(ctx)

	result, apiErr := s.postLikeRepository.LikeExists(userID, postID, dbWithCtx)
	if apiErr != nil {
		return 0, apiErr
	}
	if !result {
		return 0, &delivery.APIError{Code: constants.CreateError, Message: "user not liked this post"}
	}

	tx := dbWithCtx.Begin()
	if tx.Error != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "error starting transaction"}
	}

	likes, apiErr := s.postRepository.DecrementLikes(postID, tx)
	if apiErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, apiErr
	}

	if apiErr := s.postLikeRepository.DeleteLike(userID, postID, tx); apiErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, apiErr
	}

	if tx.Commit().Error != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return likes, nil
}
