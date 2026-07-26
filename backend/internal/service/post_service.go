package service

import (
	"context"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/model"
	"myapp/internal/repository"
	"slices"

	"gorm.io/gorm"
)

type PostServiceInterface interface {
	CreatePost(postDTO dto.PostDTO, creatorID uint, ctx context.Context) (uint, *delivery.APIError)
	GetPostByID(postID uint, ctx context.Context) (*dto.PostDTO, *delivery.APIError)
	GetCurrentUserPosts(userID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
	GetUserPostsByUsername(username string, currentUserID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
	DeletePost(postID, userID uint, ctx context.Context) *delivery.APIError
	GetPostFeed(currentUserID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
}

type PostService struct {
	postRepository     repository.PostRepositoryInterface
	postLikeRepository repository.PostLikeRepositoryInterface
	commentRepository  repository.CommentRepositoryInterface
	db                 *gorm.DB
}

func NewPostService(postRepository repository.PostRepositoryInterface,
	postLikeRepository repository.PostLikeRepositoryInterface, commentRepository repository.CommentRepositoryInterface,
	db *gorm.DB) PostServiceInterface {
	return PostService{
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		commentRepository:  commentRepository,
		db:                 db,
	}
}

func (s PostService) CreatePost(postDTO dto.PostDTO, creatorID uint, ctx context.Context) (uint, *delivery.APIError) {
	postDTO.CreatorID = creatorID
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

func (s PostService) GetCurrentUserPosts(userID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError) {
	posts, err := s.postRepository.GetAllByID(userID, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	likedPostsID, apiErr := s.postLikeRepository.GetLikedPostsID(userID, s.db.WithContext(ctx))
	if apiErr != nil {
		return nil, apiErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) GetUserPostsByUsername(username string, currentUserID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError) {
	posts, err := s.postRepository.GetAllByUsername(username, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	likedPostsID, apiErr := s.postLikeRepository.GetLikedPostsID(currentUserID, s.db.WithContext(ctx))
	if apiErr != nil {
		return nil, apiErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) DeletePost(postID, userID uint, ctx context.Context) *delivery.APIError {
	if apiErr := s.postRepository.DeletePost(postID, userID, s.db.WithContext(ctx)); apiErr != nil {
		return apiErr
	}

	return nil
}

func (s PostService) GetPostFeed(currentUserID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError) {
	posts, apiErr := s.postRepository.GetPostFeed(s.db.WithContext(ctx))
	if apiErr != nil {
		return nil, apiErr
	}

	likedPostsID, apiErr := s.postLikeRepository.GetLikedPostsID(currentUserID, s.db.WithContext(ctx))
	if apiErr != nil {
		return nil, apiErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) makePostDTOs(posts []model.Post, likedPostsID []uint) []dto.PostDTO {
	dtos := make([]dto.PostDTO, len(posts))
	for i, post := range posts {
		dtos[i] = post.ToPostDTO()
		if slices.Contains(likedPostsID, dtos[i].PostID) {
			dtos[i].IsLiked = true
		}
		commentsCount, apiErr := s.commentRepository.CountCommentsOnPost(dtos[i].PostID, s.db) //TODO: remove - sql query in for
		if apiErr == nil {
			dtos[i].CommentsCount = commentsCount
		}
	}

	return dtos
}
