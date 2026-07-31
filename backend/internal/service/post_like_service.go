package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type LikeUpdater interface {
	IncrementLikes(postID uint, db *gorm.DB) (int, *domain.DomainError)
	DecrementLikes(postID uint, db *gorm.DB) (int, *domain.DomainError)
}

type PostLikeRepository interface {
	CreateLike(likedUserID uint, postID uint, db *gorm.DB) *domain.DomainError
	DeleteLike(likedUserID, postID uint, db *gorm.DB) *domain.DomainError
	LikeExists(userID, postID uint, db *gorm.DB) (bool, *domain.DomainError)
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

func (s PostLikeService) LikePost(postID, userID uint, ctx context.Context) (int, *domain.DomainError) {
	dbWithCtx := s.db.WithContext(ctx)

	result, domainErr := s.postLikeRepository.LikeExists(userID, postID, dbWithCtx)
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

	likes, domainErr := s.postRepository.IncrementLikes(postID, tx)
	if domainErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, domainErr
	}

	if domainErr := s.postLikeRepository.CreateLike(userID, postID, tx); domainErr != nil {
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

func (s PostLikeService) DislikePost(postID, userID uint, ctx context.Context) (int, *domain.DomainError) {
	dbWithCtx := s.db.WithContext(ctx)

	result, domainErr := s.postLikeRepository.LikeExists(userID, postID, dbWithCtx)
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

	likes, domainErr := s.postRepository.DecrementLikes(postID, tx)
	if domainErr != nil {
		if result := tx.Rollback(); result.Error != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return 0, domainErr
	}

	if domainErr := s.postLikeRepository.DeleteLike(userID, postID, tx); domainErr != nil {
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
