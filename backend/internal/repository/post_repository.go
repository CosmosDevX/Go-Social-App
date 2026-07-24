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
	CreatePost(postDTO dto.PostDTO, ctx context.Context) (uint, *delivery.APIError)
	GetPostByID(postID uint, ctx context.Context) (*model.Post, *delivery.APIError)
	IncrementLikes(postID uint, ctx context.Context) (int, *delivery.APIError)
	DecrementLikes(postID uint, ctx context.Context) (int, *delivery.APIError)
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepositoryInterface {
	return PostRepository{
		db: db,
	}
}

func (r PostRepository) CreatePost(postDTO dto.PostDTO, ctx context.Context) (uint, *delivery.APIError) {
	post := model.Post{
		PostName:        postDTO.PostName,
		PostDescription: postDTO.PostDescription,
		CreatorID:       postDTO.CreatorID,
	}

	result := r.db.WithContext(ctx).Create(&post)
	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &delivery.APIError{Code: constants.CreateError, Message: "error during post creating"}
	}

	return post.ID, nil
}

func (r PostRepository) GetPostByID(postID uint, ctx context.Context) (*model.Post, *delivery.APIError) {
	var post model.Post

	result := r.db.WithContext(ctx).First(&post, "id = ?", postID)
	if result.Error != nil {
		if result.Error != nil {
			if errors.Is(result.Error, context.DeadlineExceeded) {
				return nil, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
			}
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, &delivery.APIError{Code: constants.NotFound, Message: "post not found"}
			}

			return nil, &delivery.APIError{Code: constants.FindError, Message: "error during finding post by id"}
		}
	}

	return &post, nil
}

func (r PostRepository) IncrementLikes(postID uint, ctx context.Context) (int, *delivery.APIError) {
	var likes int
	result := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes + ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if result.Error != nil {
			if errors.Is(result.Error, context.DeadlineExceeded) {
				return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
			}

			return 0, &delivery.APIError{Code: constants.UpdateError, Message: "error during increment likes on post"}
		}
	}

	if result.RowsAffected == 0 {
		return 0, &delivery.APIError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}

func (r PostRepository) DecrementLikes(postID uint, ctx context.Context) (int, *delivery.APIError) {
	var likes int
	result := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("id = ?", postID).
		Update("likes", gorm.Expr("likes - ?", 1)).
		Select("likes").
		Scan(&likes)

	if result.Error != nil {
		if result.Error != nil {
			if errors.Is(result.Error, context.DeadlineExceeded) {
				return 0, &delivery.APIError{Code: constants.RequestTimeout, Message: "request timeout"}
			}

			return 0, &delivery.APIError{Code: constants.UpdateError, Message: "error during decrement likes on post"}
		}
	}

	if result.RowsAffected == 0 {
		return 0, &delivery.APIError{Code: constants.FindError, Message: "post not found"}
	}

	return likes, nil
}
