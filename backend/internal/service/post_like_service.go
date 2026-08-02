package service

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/repository"
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
	unitOfWork         repository.UnitOfWork
	postRepository     LikeUpdater
	postLikeRepository PostLikeRepository
}

func NewPostLikeService(unitOfWork repository.UnitOfWork, postRepository LikeUpdater, postLikeRepository PostLikeRepository) PostLikeService {
	return PostLikeService{
		unitOfWork:         unitOfWork,
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
	}
}

func (s PostLikeService) LikePost(postID, userID int, ctx context.Context) (int, *domain.DomainError) {
	value, domainErr := s.unitOfWork.Do(ctx, func(ctx context.Context, repos repository.Repositories) (any, *domain.DomainError) {
		result, domainErr := repos.PostLikeRepository.LikeExists(ctx, userID, postID)
		if domainErr != nil {
			return 0, domainErr
		}
		if result {
			return 0, &domain.DomainError{Code: constants.CreateError, Message: "user already liked this post"}
		}

		likes, domainErr := repos.PostRepository.IncrementLikes(ctx, postID)
		if domainErr != nil {
			return 0, domainErr
		}

		if domainErr := repos.PostLikeRepository.CreateLike(ctx, userID, postID); domainErr != nil {
			return 0, domainErr
		}

		return likes, nil
	})

	if domainErr != nil {
		return 0, domainErr
	}

	id := value.(int)
	return id, nil
}

func (s PostLikeService) DislikePost(postID, userID int, ctx context.Context) (int, *domain.DomainError) {
	value, domainErr := s.unitOfWork.Do(ctx, func(ctx context.Context, repos repository.Repositories) (any, *domain.DomainError) {
		result, domainErr := repos.PostLikeRepository.LikeExists(ctx, userID, postID)
		if domainErr != nil {
			return 0, domainErr
		}
		if !result {
			return 0, &domain.DomainError{Code: constants.CreateError, Message: "user not liked this post"}
		}

		likes, domainErr := repos.PostRepository.DecrementLikes(ctx, postID)
		if domainErr != nil {
			return 0, domainErr
		}

		if domainErr := repos.PostLikeRepository.DeleteLike(ctx, userID, postID); domainErr != nil {
			return 0, domainErr
		}

		return likes, nil
	})

	if domainErr != nil {
		return 0, domainErr
	}

	id := value.(int)
	return id, nil
}
