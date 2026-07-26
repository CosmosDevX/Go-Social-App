package service

import (
	"context"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/model"
	"myapp/internal/repository"

	"gorm.io/gorm"
)

type CommentServiceInterface interface {
	CreateComment(commentDTO dto.CommentDTO, ctx context.Context) (uint, *delivery.APIError)
	GetAllCommentsByPostID(postID uint, ctx context.Context) ([]dto.CommentDTO, *delivery.APIError)
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

func (s CommentService) CreateComment(commentDTO dto.CommentDTO, ctx context.Context) (uint, *delivery.APIError) {
	commentID, apiErr := s.commentRepository.Create(commentDTO, s.db.WithContext(ctx))
	if apiErr != nil {
		return 0, apiErr
	}

	return commentID, nil
}

func (s CommentService) GetAllCommentsByPostID(postID uint, ctx context.Context) ([]dto.CommentDTO, *delivery.APIError) {
	comments, apiErr := s.commentRepository.GetAllByPostID(postID, s.db.WithContext(ctx))
	if apiErr != nil {
		return nil, apiErr
	}

	return s.makeCommentDTOs(comments), nil
}

func (s CommentService) makeCommentDTOs(comments []model.Comment) []dto.CommentDTO {
	dtos := make([]dto.CommentDTO, len(comments))
	for i, comment := range comments {
		dtos[i] = comment.ToCommentDTO()
	}

	return dtos
}
