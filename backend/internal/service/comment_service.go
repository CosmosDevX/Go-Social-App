package service

import (
	"context"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
)

type CommentRepository interface {
	Delete(ctx context.Context, commentID, userID int) *domain.DomainError
	Create(ctx context.Context, commentDTO dto.CommentDTO) (int, *domain.DomainError)
	GetAllByPostID(ctx context.Context, postID int) ([]domain.Comment, *domain.DomainError)
}

type UsernamesGetter interface {
	GetUsernameByIDs(ctx context.Context, ids []int) ([]string, *domain.DomainError)
}

type CommentService struct {
	commentRepository CommentRepository
	usernamesGetter   UsernamesGetter
}

func NewCommentService(commentRepository CommentRepository, usernamesGetter UsernamesGetter) CommentService {
	return CommentService{
		commentRepository: commentRepository,
		usernamesGetter:   usernamesGetter,
	}
}

func (s CommentService) CreateComment(commentDTO dto.CommentDTO, creatorID, postID int, ctx context.Context) (int, *domain.DomainError) {
	commentDTO.CreatorID = creatorID
	commentDTO.PostID = postID
	commentID, domainErr := s.commentRepository.Create(ctx, commentDTO)
	if domainErr != nil {
		return 0, domainErr
	}

	return commentID, nil
}

func (s CommentService) DeleteComment(commentID, userID int, ctx context.Context) *domain.DomainError {
	if domainErr := s.commentRepository.Delete(ctx, commentID, userID); domainErr != nil {
		return domainErr
	}

	return nil
}

func (s CommentService) GetAllCommentsByPostID(postID int, ctx context.Context) ([]dto.CommentDTO, *domain.DomainError) {
	comments, domainErr := s.commentRepository.GetAllByPostID(ctx, postID)
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
