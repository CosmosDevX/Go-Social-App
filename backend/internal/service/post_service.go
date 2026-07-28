package service

import (
	"context"
	"log"
	"mime/multipart"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/model"
	"myapp/internal/repository"
	"myapp/internal/utils"
	"slices"

	"gorm.io/gorm"
)

type PostServiceInterface interface {
	CreatePost(postDTO dto.PostDTO, creatorID uint, file multipart.File, header *multipart.FileHeader, ctx context.Context) (uint, *delivery.APIError)
	GetPostByID(postID uint, ctx context.Context) (*dto.PostDTO, *delivery.APIError)
	GetCurrentUserPosts(userID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
	GetUserPostsByUsername(username string, currentUserID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
	DeletePost(postID, userID uint, ctx context.Context) *delivery.APIError
	GetPostFeed(currentUserID uint, ctx context.Context) ([]dto.PostDTO, *delivery.APIError)
}

type PostService struct {
	fileManager        utils.FileManagerInterface
	postRepository     repository.PostRepositoryInterface
	postLikeRepository repository.PostLikeRepositoryInterface
	commentRepository  repository.CommentRepositoryInterface
	db                 *gorm.DB
}

func NewPostService(postRepository repository.PostRepositoryInterface,
	postLikeRepository repository.PostLikeRepositoryInterface, commentRepository repository.CommentRepositoryInterface,
	fileManager utils.FileManagerInterface, db *gorm.DB) PostServiceInterface {
	return PostService{
		fileManager:        fileManager,
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		commentRepository:  commentRepository,
		db:                 db,
	}
}

func (s PostService) CreatePost(postDTO dto.PostDTO, creatorID uint, file multipart.File, header *multipart.FileHeader, ctx context.Context) (uint, *delivery.APIError) {
	tx := s.db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "error during start transaction"}
	}

	filename, err := s.fileManager.SaveFile(file, header, "uploads")
	if err != nil {
		return 0, &delivery.APIError{Code: constants.SaveError, Message: err.Error()}
	}

	postDTO.CreatorID = creatorID
	postDTO.ImageName = filename
	postID, apiErr := s.postRepository.Create(postDTO, tx)
	if apiErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}

		if err := s.fileManager.DeleteFile("/uploads", filename); err != nil {
			log.Println(err.Error())
		}

		return 0, apiErr
	}

	if err := tx.Commit().Error; err != nil {
		return 0, &delivery.APIError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return postID, nil
}

func (s PostService) DeletePost(postID, userID uint, ctx context.Context) *delivery.APIError {
	tx := s.db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return &delivery.APIError{Code: constants.TransactionError, Message: "error during start transaction"}
	}

	imageName, apiErr := s.postRepository.GetImageName(postID, tx)
	if apiErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return apiErr
	}

	if apiErr := s.postRepository.DeletePost(postID, userID, tx); apiErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return apiErr
	}

	if err := s.fileManager.DeleteFile("/uploads", imageName); err != nil {
		if err := tx.Rollback().Error; err != nil {
			return &delivery.APIError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return &delivery.APIError{Code: constants.DeleteError, Message: "post image not deleted"}
	}

	if err := tx.Commit().Error; err != nil {
		return &delivery.APIError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return nil
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
