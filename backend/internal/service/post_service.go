package service

import (
	"context"
	"log"
	"mime/multipart"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/internal/utils"
	"slices"
)

type CommentCounter interface {
	CountCommentsOnPost(ctx context.Context, postID int) (int, *domain.DomainError)
}

type PostRepository interface {
	Create(ctx context.Context, postDTO dto.PostDTO) (int, *domain.DomainError)
	GetByID(ctx context.Context, postID int) (*domain.Post, *domain.DomainError)
	GetAllByID(ctx context.Context, userID int) ([]domain.Post, *domain.DomainError)
	GetAllByUsername(ctx context.Context, username string, userID int) ([]domain.Post, *domain.DomainError)
	DeletePost(ctx context.Context, postID, userID int) *domain.DomainError
	GetPostFeed(ctx context.Context) ([]domain.Post, *domain.DomainError)
	GetImageName(ctx context.Context, postID int) (string, *domain.DomainError)
}

type PostLikeGetter interface {
	GetLikedPostsID(ctx context.Context, userID int) ([]int, *domain.DomainError)
}

type UsernameGetter interface {
	GetUsernameByID(ctx context.Context, userID int) (string, *domain.DomainError)
}

type PostService struct {
	unitOfWork         repository.UnitOfWork
	fileManager        utils.FileManagerInterface
	postRepository     PostRepository
	postLikeRepository PostLikeGetter
	commentRepository  CommentCounter
	usernameGetter     UsernameGetter
}

func NewPostService(unitOfWork repository.UnitOfWork, postRepository PostRepository,
	postLikeRepository PostLikeGetter, commentRepository CommentCounter,
	fileManager utils.FileManagerInterface, usernameGetter UsernameGetter) PostService {
	return PostService{
		unitOfWork:         unitOfWork,
		fileManager:        fileManager,
		postRepository:     postRepository,
		postLikeRepository: postLikeRepository,
		commentRepository:  commentRepository,
		usernameGetter:     usernameGetter,
	}
}

func (s PostService) CreatePost(ctx context.Context, postDTO dto.PostDTO, creatorID int, file multipart.File, header *multipart.FileHeader) (int, *domain.DomainError) {
	value, domainErr := s.unitOfWork.Do(ctx, func(ctx context.Context, repos repository.Repositories) (any, *domain.DomainError) {
		filename, err := s.fileManager.SaveFile(file, header, "uploads")
		if err != nil {
			return 0, &domain.DomainError{Code: constants.SaveError, Message: err.Error()}
		}

		postDTO.CreatorID = creatorID
		postDTO.ImageName = filename
		postID, domainErr := s.postRepository.Create(ctx, postDTO)
		if domainErr != nil {
			if err := s.fileManager.DeleteFile("/uploads", filename); err != nil {
				log.Println(err.Error())
			}

			return 0, domainErr
		}

		return postID, nil
	})

	if domainErr != nil {
		return 0, domainErr
	}

	postID := value.(int)
	return postID, nil
}

func (s PostService) DeletePost(ctx context.Context, postID, userID int) *domain.DomainError {
	_, domainErr := s.unitOfWork.Do(ctx, func(ctx context.Context, repos repository.Repositories) (any, *domain.DomainError) {
		imageName, domainErr := s.postRepository.GetImageName(ctx, postID)
		if domainErr != nil {
			return nil, domainErr
		}

		if domainErr := s.postRepository.DeletePost(ctx, postID, userID); domainErr != nil {
			return nil, domainErr
		}

		if err := s.fileManager.DeleteFile("/uploads", imageName); err != nil {
			return nil, &domain.DomainError{Code: constants.DeleteError, Message: "post image not deleted"}
		}

		return nil, nil
	})

	if domainErr != nil {
		return domainErr
	}

	return nil
}

func (s PostService) GetPostByID(ctx context.Context, postID int) (*dto.PostDTO, *domain.DomainError) {
	post, err := s.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	dto := post.ToPostDTO()
	return &dto, nil
}

func (s PostService) GetCurrentUserPosts(ctx context.Context, userID int) ([]dto.PostDTO, *domain.DomainError) {
	posts, domainErr := s.postRepository.GetAllByID(ctx, userID)
	if domainErr != nil {
		return nil, domainErr
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(ctx, userID)
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) GetUserPostsByUsername(ctx context.Context, username string, currentUserID int) ([]dto.PostDTO, *domain.DomainError) {
	posts, err := s.postRepository.GetAllByUsername(ctx, username, currentUserID)
	if err != nil {
		return nil, err
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(ctx, currentUserID)
	if domainErr != nil {
		return nil, domainErr
	}

	return s.makePostDTOs(posts, likedPostsID), nil
}

func (s PostService) GetPostFeed(ctx context.Context, currentUserID int) ([]dto.PostDTO, *domain.DomainError) {
	posts, domainErr := s.postRepository.GetPostFeed(ctx)
	if domainErr != nil {
		return nil, domainErr
	}

	likedPostsID, domainErr := s.postLikeRepository.GetLikedPostsID(ctx, currentUserID)
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
