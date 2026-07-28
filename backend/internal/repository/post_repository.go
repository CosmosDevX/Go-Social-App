package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type PostRepositoryInterface interface {
	Create(postDTO dto.PostDTO, db *gorm.DB) (uint, *delivery.APIError)
	GetByID(postID uint, db *gorm.DB) (*model.Post, *delivery.APIError)
	GetAllByID(userID uint, db *gorm.DB) ([]model.Post, *delivery.APIError)
	GetAllByUsername(username string, db *gorm.DB) ([]model.Post, *delivery.APIError)
	IncrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError)
	DecrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError)
	DeletePost(postID, userID uint, db *gorm.DB) *delivery.APIError
	GetPostFeed(db *gorm.DB) ([]model.Post, *delivery.APIError)
	GetImageName(postID uint, db *gorm.DB) (string, *delivery.APIError)
}

type PostRepository struct{}

func (r PostRepository) Create(postDTO dto.PostDTO, db *gorm.DB) (uint, *delivery.APIError) {
	post := model.Post{
		PostName:        postDTO.PostName,
		PostDescription: postDTO.PostDescription,
		CreatorID:       postDTO.CreatorID,
		ImageName:       postDTO.ImageName,
	}

	result := db.Create(&post)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.CreateError, Message: "error during post creating"}
	}

	return post.ID, nil
}

func (r PostRepository) GetByID(postID uint, db *gorm.DB) (*model.Post, *delivery.APIError) {
	var post model.Post

	result := db.First(&post, "id = ?", postID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "post not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding post by id"}
	}

	return &post, nil
}

func (r PostRepository) GetAllByID(userID uint, db *gorm.DB) ([]model.Post, *delivery.APIError) {
	var posts []model.Post

	result := db.Preload("Creator").Find(&posts, "creator_id = ?", userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "posts not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding posts by user id"}
	}

	return posts, nil
}

func (r PostRepository) GetAllByUsername(username string, db *gorm.DB) ([]model.Post, *delivery.APIError) {
	var posts []model.Post

	result := db.Preload("Creator").Joins("JOIN users ON users.id = posts.creator_id").
		Where("users.username = ?", username).
		Find(&posts)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &delivery.APIError{Code: constants.NotFound, Message: "posts not found"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding posts by username"}
	}

	return posts, nil
}

func (r PostRepository) IncrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError) {
	var likes int

	result := db.Model(&model.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes + ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.UpdateError, Message: "error during increment likes on post"}
	}

	if result.RowsAffected == 0 {
		return 0, &delivery.APIError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}

func (r PostRepository) DecrementLikes(postID uint, db *gorm.DB) (int, *delivery.APIError) {
	var likes int
	result := db.Model(&model.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes - ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.UpdateError, Message: "error during decrement likes on post"}
	}

	if result.RowsAffected == 0 {
		return 0, &delivery.APIError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}

func (r PostRepository) DeletePost(postID, userID uint, db *gorm.DB) *delivery.APIError {
	result := db.Delete(&model.Post{}, "id = ? AND creator_id = ?", postID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &delivery.APIError{Code: constants.DeleteError, Message: "error during deleting the post"}
	}

	if result.RowsAffected == 0 {
		return &delivery.APIError{Code: constants.NotFound, Message: "post not deleted"}
	}

	return nil
}

func (r PostRepository) GetPostFeed(db *gorm.DB) ([]model.Post, *delivery.APIError) {
	var posts []model.Post
	result := db.
		Preload("Creator").
		Order("RANDOM()").
		Limit(30).
		Find(&posts)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return nil, &delivery.APIError{Code: constants.FindError, Message: "error during getting post feed"}
	}

	return posts, nil
}

func (r PostRepository) GetImageName(postID uint, db *gorm.DB) (string, *delivery.APIError) {
	var imageName string
	result := db.Model(&model.Post{}).Where("id = ?", postID).Select("image_name").Scan(&imageName)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return "", &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", &delivery.APIError{Code: constants.NotFound, Message: "post not found"}
		}

		return "", &delivery.APIError{Code: constants.FindError, Message: "error during getting post image"}
	}

	return imageName, nil
}
