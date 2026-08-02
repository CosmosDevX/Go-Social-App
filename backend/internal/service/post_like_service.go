package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type LikeUpdater interface {
	IncrementLikes(ctx context.Context, postID int) (int, *domain.DomainError)
	DecrementLikes(ctx context.Context, postID int) (int, *domain.DomainError)
}

type PostLikeRepository interface {
	CreateLike(ctx context.Context, likedUserID, postID int) *domain.DomainError
	DeleteLike(ctx context.Context, likedUserID, postID int) *domain.DomainError
	LikeExists(ctx context.Context, userID, postID int) (bool, *domain.DomainError)
}

type PostLikeService struct {
	postRepository     LikeUpdater
	postLikeRepository PostLikeRepository
	db                 *gorm.DB
}

func NewPostLikeService(postRepository LikeUpdater, postLikeRepository PostLikeRepository, db *gorm.DB) PostLikeService {
	return PostLikeService{
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		db:                 db,
	}
}

func (s PostLikeService) LikePost(postID, userID int, ctx context.Context) (int, *domain.DomainError) {
	dbWithCtx := s.db.WithContext(ctx)

	result, domainErr := s.postLikeRepository.LikeExists(ctx, userID, postID)
	if domainErr != nil {
		return 0, domainErr
	}
	if result {
		return 0, &domain.DomainError{Code: constants.CreateError, Message: "user already liked this post"}
	}

	tx := dbWithCtx.Begin()
	if tx.Error != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "error starting transaction"}
	}

	likes, domainErr := s.postRepository.IncrementLikes(ctx, postID)
	if domainErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, domainErr
	}

	if domainErr := s.postLikeRepository.CreateLike(ctx, userID, postID); domainErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, domainErr
	}

	if tx.Commit().Error != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return likes, nil
}

func (s PostLikeService) DislikePost(postID, userID int, ctx context.Context) (int, *domain.DomainError) {
	dbWithCtx := s.db.WithContext(ctx)

	result, domainErr := s.postLikeRepository.LikeExists(ctx, userID, postID)
	if domainErr != nil {
		return 0, domainErr
	}
	if !result {
		return 0, &domain.DomainError{Code: constants.CreateError, Message: "user not liked this post"}
	}

	tx := dbWithCtx.Begin()
	if tx.Error != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "error starting transaction"}
	}

	likes, domainErr := s.postRepository.DecrementLikes(ctx, postID)
	if domainErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, domainErr
	}

	if domainErr := s.postLikeRepository.DeleteLike(ctx, userID, postID); domainErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, domainErr
	}

	if tx.Commit().Error != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return likes, nil
}
