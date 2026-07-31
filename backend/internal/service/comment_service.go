package service

import (
	"context"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Delete(commentID, userID uint, db *gorm.DB) *domain.DomainError
	Create(commentDTO dto.CommentDTO, db *gorm.DB) (uint, *domain.DomainError)
	GetAllByPostID(postID uint, db *gorm.DB) ([]domain.Comment, *domain.DomainError)
}

type CommentService struct {
	commentRepository CommentRepository
	db                *gorm.DB
}

func NewCommentService(commentRepository CommentRepository, db *gorm.DB) CommentService {
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
