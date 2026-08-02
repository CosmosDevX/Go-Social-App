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
	CountCommentsOnPost(ctx context.Context, postID int) (int, *domain.DomainError)
}

type PostRepository interface {
	Create(ctx context.Context, postDTO dto.PostDTO) (int, *domain.DomainError)
	GetByID(ctx context.Context, postID int) (*domain.Post, *domain.DomainError)
	GetAllByID(ctx context.Context, userID int) ([]domain.Post, *domain.DomainError)
	GetAllByUsername(ctx context.Context, username string) ([]domain.Post, *domain.DomainError)
	DeletePost(ctx context.Context, postID, userID int) *domain.DomainError
	GetPostFeed(db *gorm.DB) ([]domain.Post, *domain.DomainError)
	GetImageName(ctx context.Context, postID int) (string, *domain.DomainError)
}

type PostLikeGetter interface {
	GetLikedPostsID(ctx context.Context, userID int) ([]int, *domain.DomainError)
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

func (s PostService) CreatePost(postDTO dto.PostDTO, creatorID int, file multipart.File, header *multipart.FileHeader, ctx context.Context) (int, *domain.DomainError) {
	tx := s.db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return 0, &domain.DomainError{Code: constants.TransactionError, Message: "error during start transaction"}
	}

	filename, err := s.fileManager.SaveFile(file, header, "uploads")
	if err != nil {
		return 0, &domain.DomainError{Code: constants.SaveError, Message: err.Error()}
	}

	postDTO.CreatorID = int(creatorID)
	postDTO.ImageName = filename
	postID, domainErr := s.postRepository.Create(ctx, postDTO)
	if domainErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return 0, &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}

		log.Println("file delete and transaction rollback")
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

func (s PostService) DeletePost(postID, userID int, ctx context.Context) *domain.DomainError {
	tx := s.db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return &domain.DomainError{Code: constants.TransactionError, Message: "error during start transaction"}
	}

	imageName, domainErr := s.postRepository.GetImageName(ctx, postID)
	if domainErr != nil {
		if err := tx.Rollback().Error; err != nil {
			return &domain.DomainError{Code: constants.TransactionError, Message: "transaction rollback failed"}
		}
		return domainErr
	}

	if domainErr := s.postRepository.DeletePost(ctx, postID, userID); domainErr != nil {
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

func (s PostService) GetPostByID(postID int, ctx context.Context) (*dto.PostDTO, *domain.DomainError) {
	post, err := s.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	dto := post.ToPostDTO()
	return &dto, nil
}

func (s PostService) GetCurrentUserPosts(userID int, ctx context.Context) ([]dto.PostDTO, *domain.DomainError) {
	posts, err := s.postRepository.GetAllByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(ctx, userID)
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) GetUserPostsByUsername(username string, currentUserID int, ctx context.Context) ([]dto.PostDTO, *domain.DomainError) {
	posts, err := s.postRepository.GetAllByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(ctx, currentUserID)
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

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(ctx, int(currentUserID))
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) makePostDTOs(posts []domain.Post, likedPostsID []int) []dto.PostDTO {
	dtos := make([]dto.PostDTO, len(posts))
	for i, post := range posts {
		dtos[i] = post.ToPostDTO()
		if slices.Contains(likedPostsID, dtos[i].PostID) {
			dtos[i].IsLiked = true
		}
		commentsCount, domainErr := s.commentRepository.CountCommentsOnPost(context.TODO(), int(dtos[i].PostID)) //TODO: remove - sql query in for
		if domainErr == nil {
			dtos[i].CommentsCount = commentsCount
		}
	}

	return dtos
}
