package service

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/repository"
	"strconv"

	"gorm.io/gorm"
)

type PostServiceInterface interface {
	CreatePost(postDTO dto.PostDTO, ctx context.Context) (uint, *delivery.APIError)
	GetPostByID(postID uint, ctx context.Context) (*dto.PostDTO, *delivery.APIError)
	GetAllUserPosts(ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
	LikePost(postID uint, ctx context.Context) (int, *delivery.APIError)
}

type PostService struct {
	postRepository     repository.PostRepositoryInterface
	postLikeRepository repository.PostLikeRepositoryInterface
	db                 *gorm.DB
}

func NewPostService(postRepository repository.PostRepositoryInterface, postLikeRepository repository.PostLikeRepositoryInterface, db *gorm.DB) PostServiceInterface {
	return PostService{
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		db:                 db,
	}
}

func (s PostService) CreatePost(postDTO dto.PostDTO, ctx context.Context) (uint, *delivery.APIError) {
	parsedUserID, parseErr := s.parseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		return 0, &delivery.APIError{Code: constants.ParseError, Message: parseErr.Error()}
	}

	postDTO.CreatorID = parsedUserID
	postID, apiErr := s.postRepository.Create(postDTO, s.db.WithContext(ctx))
	if apiErr != nil {
		return 0, apiErr
	}

	return postID, nil
}

func (s PostService) GetPostByID(postID uint, ctx context.Context) (*dto.PostDTO, *delivery.APIError) {
	post, err := s.postRepository.GetByID(postID, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	dto := post.ToPostDTO()
	return &dto, nil
}

func (s PostService) GetAllUserPosts(ctx context.Context) ([]dto.PostDTO, *delivery.APIError) {
	parsedUserID, parseErr := s.parseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		return nil, &delivery.APIError{Code: constants.ParseError, Message: parseErr.Error()}
	}

	posts, err := s.postRepository.GetAll(parsedUserID, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.PostDTO, len(posts))
	for i, post := range posts {
		dtos[i] = post.ToPostDTO()
	}

	return dtos, nil
}

func (s PostService) LikePost(postID uint, ctx context.Context) (int, *delivery.APIError) {
	parsedUserID, parseErr := s.parseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		return 0, &delivery.APIError{Code: constants.ParseError, Message: parseErr.Error()}
	}

	dbWithCtx := s.db.WithContext(ctx)

	result, apiErr := s.postLikeRepository.LikeExists(parsedUserID, postID, dbWithCtx)
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
		tx.Rollback()
		return 0, apiErr
	}

	if apiErr := s.postLikeRepository.CreateLike(parsedUserID, postID, tx); apiErr != nil {
		tx.Rollback()
		return 0, apiErr
	}

	if tx.Commit().Error != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return likes, nil
}

func (s PostService) parseUserID(ctxValue any) (uint, error) {
	stringUserID, ok := ctxValue.(string)
	if !ok {
		return 0, errors.New("error during parsing userID to string")
	}

	userID, err := strconv.ParseUint(stringUserID, 10, 64)
	if err != nil {
		return 0, errors.New("error during parsing userID to uint")
	}

	return uint(userID), nil
}
