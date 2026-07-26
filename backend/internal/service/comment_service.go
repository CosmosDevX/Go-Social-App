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
	DeleteComment(commentID, userID uint, ctx context.Context) *delivery.APIError
	CreateComment(commentDTO dto.CommentDTO, creatorID, postID uint, ctx context.Context) (uint, *delivery.APIError)
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

func (s CommentService) CreateComment(commentDTO dto.CommentDTO, creatorID, postID uint, ctx context.Context) (uint, *delivery.APIError) {
	commentDTO.CreatorID = creatorID
	commentDTO.PostID = postID
	commentID, apiErr := s.commentRepository.Create(commentDTO, s.db.WithContext(ctx))
	if apiErr != nil {
		return 0, apiErr
	}

	return commentID, nil
}

func (s CommentService) DeleteComment(commentID, userID uint, ctx context.Context) *delivery.APIError {
	if apiErr := s.commentRepository.Delete(commentID, userID, s.db.WithContext(ctx)); apiErr != nil {
		return apiErr
	}

	return nil
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
