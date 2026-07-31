package service

import (
	"context"
	"log"
	"mime/multipart"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"slices"

	"gorm.io/gorm"
)

type CommentCounter interface {
	CountCommentsOnPost(postID uint, db *gorm.DB) (int, *domain.DomainError)
}

type PostRepository interface {
	Create(postDTO dto.PostDTO, db *gorm.DB) (uint, *domain.DomainError)
	GetByID(postID uint, db *gorm.DB) (*domain.Post, *domain.DomainError)
	GetAllByID(userID uint, db *gorm.DB) ([]domain.Post, *domain.DomainError)
	GetAllByUsername(username string, db *gorm.DB) ([]domain.Post, *domain.DomainError)
	DeletePost(postID, userID uint, db *gorm.DB) *domain.DomainError
	GetPostFeed(db *gorm.DB) ([]domain.Post, *domain.DomainError)
	GetImageName(postID uint, db *gorm.DB) (string, *domain.DomainError)
}

type PostLikeGetter interface {
	GetLikedPostsID(userID uint, db *gorm.DB) ([]uint, *domain.DomainError)
}

type PostService struct {
	fileManager        utils.FileManagerInterface
	postRepository     PostRepository
	postLikeRepository PostLikeGetter
	commentRepository  CommentCounter
	db                 *gorm.DB
}

func NewPostService(postRepository PostRepository,
	postLikeRepository PostLikeGetter, commentRepository CommentCounter,
	fileManager utils.FileManagerInterface, db *gorm.DB) PostService {
	return PostService{
		fileManager:        fileManager,
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		commentRepository:  commentRepository,
		db:                 db,
	}
}

func (s PostService) CreatePost(postDTO dto.PostDTO, creatorID uint, file multipart.File, header *multipart.FileHeader, ctx context.Context) (uint, *domain.DomainError) {
	tx := s.db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "error during start transaction"}
	}

	filename, err := s.fileManager.SaveFile(file, header, "uploads")
	if err != nil {
		return 0, &domain.DomainError{Code: constants.SaveError, Message: err.Error()}
	}

	postDTO.CreatorID = creatorID
	postDTO.ImageName = filename
	postID, domainErr := s.postRepository.Create(postDTO, tx)
	if domainErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}

		if err := s.fileManager.DeleteFile("/uploads", filename); err != nil {
			log.Println(err.Error())
		}

		return 0, domainErr
	}

	if err := tx.Commit().Error; err != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return postID, nil
}

func (s PostService) DeletePost(postID, userID uint, ctx context.Context) *domain.DomainError {
	tx := s.db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return &domain.DomainError{Code: constants.TransactionError, Message: "error during start transaction"}
	}

	imageName, domainErr := s.postRepository.GetImageName(postID, tx)
	if domainErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return domainErr
	}

	if domainErr := s.postRepository.DeletePost(postID, userID, tx); domainErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return domainErr
	}

	if err := s.fileManager.DeleteFile("/uploads", imageName); err != nil {
		if err := tx.Rollback().Error; err != nil {
			return &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return &domain.DomainError{Code: constants.DeleteError, Message: "post image not deleted"}
	}

	if err := tx.Commit().Error; err != nil {
		return &domain.DomainError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return nil
}

func (s PostService) GetPostByID(postID uint, ctx context.Context) (*dto.PostDTO, *domain.DomainError) {
	post, err := s.postRepository.GetByID(postID, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	dto := post.ToPostDTO()
	return &dto, nil
}

func (s PostService) GetCurrentUserPosts(userID uint, ctx context.Context) ([]dto.PostDTO, *domain.DomainError) {
	posts, err := s.postRepository.GetAllByID(userID, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(userID, s.db.WithContext(ctx))
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) GetUserPostsByUsername(username string, currentUserID uint, ctx context.Context) ([]dto.PostDTO, *domain.DomainError) {
	posts, err := s.postRepository.GetAllByUsername(username, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(currentUserID, s.db.WithContext(ctx))
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) GetPostFeed(currentUserID uint, ctx context.Context) ([]dto.PostDTO, *domain.DomainError) {
	posts, domainErr := s.postRepository.GetPostFeed(s.db.WithContext(ctx))
	if domainErr != nil {
		return nil, domainErr
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(currentUserID, s.db.WithContext(ctx))
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) makePostDTOs(posts []domain.Post, likedPostsID []uint) []dto.PostDTO {
	dtos := make([]dto.PostDTO, len(posts))
	for i, post := range posts {
		dtos[i] = post.ToPostDTO()
		if slices.Contains(likedPostsID, dtos[i].PostID) {
			dtos[i].IsLiked = true
		}
		commentsCount, domainErr := s.commentRepository.CountCommentsOnPost(dtos[i].PostID, s.db) //TODO: remove - sql query in for
		if domainErr == nil {
			dtos[i].CommentsCount = commentsCount
		}
	}

	return dtos
}
