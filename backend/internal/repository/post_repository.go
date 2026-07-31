package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"gorm.io/gorm"
)

type PostRepository struct{}

func (r PostRepository) Create(postDTO dto.PostDTO, db *gorm.DB) (uint, *domain.DomainError) {
	post := domain.Post{
		PostName:        postDTO.PostName,
		PostDescription: postDTO.PostDescription,
		CreatorID:       postDTO.CreatorID,
		ImageName:       postDTO.ImageName,
	}

	result := db.Create(&post)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during post creating"}
	}

	return post.ID, nil
}

func (r PostRepository) GetByID(postID uint, db *gorm.DB) (*domain.Post, *domain.DomainError) {
	var post domain.Post

	result := db.First(&post, "id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "post not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during finding post by id"}
	}

	return &post, nil
}

func (r PostRepository) GetAllByID(userID uint, db *gorm.DB) ([]domain.Post, *domain.DomainError) {
	var posts []domain.Post

	result := db.Preload("Creator").Find(&posts, "creator_id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "posts not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during finding posts by user id"}
	}

	return posts, nil
}

func (r PostRepository) GetAllByUsername(username string, db *gorm.DB) ([]domain.Post, *domain.DomainError) {
	var posts []domain.Post

	result := db.Preload("Creator").Joins("JOIN users ON users.id = posts.creator_id").
		Where("users.username = ?", username).
		Find(&posts)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "posts not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during finding posts by username"}
	}

	return posts, nil
}

func (r PostRepository) IncrementLikes(postID uint, db *gorm.DB) (int, *domain.DomainError) {
	var likes int

	result := db.Model(&domain.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes + ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.UpdateError, Message: "error during increment likes on post"}
	}

	if result.RowsAffected == 0 {
		return 0, &domain.DomainError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}

func (r PostRepository) DecrementLikes(postID uint, db *gorm.DB) (int, *domain.DomainError) {
	var likes int
	result := db.Model(&domain.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes - ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.UpdateError, Message: "error during decrement likes on post"}
	}

	if result.RowsAffected == 0 {
		return 0, &domain.DomainError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}

func (r PostRepository) DeletePost(postID, userID uint, db *gorm.DB) *domain.DomainError {
	result := db.Delete(&domain.Post{}, "id = ? AND creator_id = ?", postID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &domain.DomainError{Code: constants.DeleteError, Message: "error during deleting the post"}
	}

	if result.RowsAffected == 0 {
		return &domain.DomainError{Code: constants.NotFound, Message: "post not deleted"}
	}

	return nil
}

func (r PostRepository) GetPostFeed(db *gorm.DB) ([]domain.Post, *domain.DomainError) {
	var posts []domain.Post
	result := db.
		Preload("Creator").
		Order("RANDOM()").
		Limit(30).
		Find(&posts)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during getting post feed"}
	}

	return posts, nil
}

func (r PostRepository) GetImageName(postID uint, db *gorm.DB) (string, *domain.DomainError) {
	var imageName string
	result := db.Model(&domain.Post{}).Where("id = ?", postID).Select("image_name").Scan(&imageName)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return "", &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", &domain.DomainError{Code: constants.NotFound, Message: "post not found"}
		}

		return "", &domain.DomainError{Code: constants.FindError, Message: "error during getting post image"}
	}

	return imageName, nil
}
