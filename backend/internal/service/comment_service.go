package service

import (
	"context"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/repository"

	"gorm.io/gorm"
)

type CommentServiceInterface interface {
	DeleteComment(commentID, userID uint, ctx context.Context) *domain.DomainError
	CreateComment(commentDTO dto.CommentDTO, creatorID, postID uint, ctx context.Context) (uint, *domain.DomainError)
	GetAllCommentsByPostID(postID uint, ctx context.Context) ([]dto.CommentDTO, *domain.DomainError)
}

type CommentService struct {
	commentRepository repository.CommentRepositoryInterface
	db                *gorm.DB
}

func NewCommentService(commentRepository repository.CommentRepositoryInterface, db *gorm.DB) CommentServiceInterface {
	return CommentService{
		commentRepository: commentRepository,
		db:                db,
	}
}

func (s CommentService) CreateComment(commentDTO dto.CommentDTO, creatorID, postID uint, ctx context.Context) (uint, *domain.DomainError) {
	commentDTO.CreatorID = creatorID
	commentDTO.PostID = postID
	commentID, domainErr := s.commentRepository.Create(commentDTO, s.db.WithContext(ctx))
	if domainErr != nil {
		return 0, domainErr
	}

	return commentID, nil
}

func (s CommentService) DeleteComment(commentID, userID uint, ctx context.Context) *domain.DomainError {
	if domainErr := s.commentRepository.Delete(commentID, userID, s.db.WithContext(ctx)); domainErr != nil {
		return domainErr
	}

	return nil
}

func (s CommentService) GetAllCommentsByPostID(postID uint, ctx context.Context) ([]dto.CommentDTO, *domain.DomainError) {
	comments, domainErr := s.commentRepository.GetAllByPostID(postID, s.db.WithContext(ctx))
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makeCommentDTOs(comments), nil
}

func (s CommentService) makeCommentDTOs(comments []domain.Comment) []dto.CommentDTO {
	dtos := make([]dto.CommentDTO, len(comments))
	for i, comment := range comments {
		dtos[i] = comment.ToCommentDTO()
	}

	return dtos
}
