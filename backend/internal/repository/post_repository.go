package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type UserFieldsGetter interface {
	GetUsernameByID(ctx context.Context, userID int) (string, *domain.DomainError)
	GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError)
}

type PostRepository struct {
	db               *sqlx.DB
	userFieldsGetter UserFieldsGetter
}

func NewPostRepository(db *sqlx.DB, userFieldsGetter UserFieldsGetter) PostRepository {
	return PostRepository{
		db:               db,
		userFieldsGetter: userFieldsGetter,
	}
}

func (r PostRepository) Create(ctx context.Context, postDTO dto.PostDTO) (int, *domain.DomainError) {
	query := `INSERT INTO posts(name, description, creator_id, image_name) VALUES($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, postDTO.PostName, postDTO.PostDescription, postDTO.CreatorID, postDTO.ImageName).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during create post"}
	}

	return id, nil
}

func (r PostRepository) GetByID(ctx context.Context, postID int) (*domain.Post, *domain.DomainError) {
	query := `SELECT * FROM posts WHERE id = $1`
	var post domain.Post
	err := r.db.GetContext(ctx, &post, query, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "post not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during get post by id"}
	}

	return &post, nil
}

func (r PostRepository) GetAllByID(ctx context.Context, userID int) ([]domain.Post, *domain.DomainError) {
	query := `SELECT * FROM posts WHERE creator_id = $1`
	var posts []domain.Post
	err := r.db.SelectContext(ctx, &posts, query, userID)
	if err != nil {
		log.Println(err)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "posts not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during fing posts by user id"}
	}

	username, domainErr := r.userFieldsGetter.GetUsernameByID(ctx, userID)
	if domainErr != nil {
		return nil, domainErr
	}

	for i := range posts {
		posts[i].CreatorUsername = username
	}

	return posts, nil
}

func (r PostRepository) GetAllByUsername(ctx context.Context, username string) ([]domain.Post, *domain.DomainError) {
	user, domainErr := r.userFieldsGetter.GetUserByName(ctx, username)
	if domainErr != nil {
		return nil, domainErr
	}
	userID := user.ID

	posts, domainErr := r.GetAllByID(ctx, userID)
	if domainErr != nil {
		return nil, domainErr
	}

	return posts, nil
}

func (r PostRepository) IncrementLikes(ctx context.Context, postID int) (int, *domain.DomainError) {
	query := `UPDATE posts SET likes = likes + 1 WHERE id = $1 RETURNING likes`
	var likes int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&likes)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.UpdateError, Message: "error during increment likes on post"}
	}

	return likes, nil
}

func (r PostRepository) DecrementLikes(ctx context.Context, postID int) (int, *domain.DomainError) {
	query := `UPDATE posts SET likes = likes - 1 WHERE id = $1 RETURNING likes`
	var likes int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&likes)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.UpdateError, Message: "error during decrement likes on post"}
	}

	return likes, nil
}

func (r PostRepository) DeletePost(ctx context.Context, postID, userID int) *domain.DomainError {
	query := `DELETE FROM posts WHERE id = $1 AND creator_id = $2`
	sqlResult, err := r.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &domain.DomainError{Code: constants.DeleteError, Message: "error during delete post"}
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return &domain.DomainError{Code: constants.DatabaseError, Message: "error during get affected rows"}
	}

	if rowsAffected == 0 {
		return &domain.DomainError{Code: constants.NotFound, Message: "post not deleted"}
	}

	return nil
}

func (r PostRepository) GetPostFeed(db *gorm.DB) ([]domain.Post, *domain.DomainError) {
	var posts []domain.Post
	result := db. //TODO: remove and make post feed pagination
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

func (r PostRepository) GetImageName(ctx context.Context, postID int) (string, *domain.DomainError) {
	query := `SELECT image_name FROM posts WHERE id = $1`
	var post domain.Post
	err := r.db.GetContext(ctx, &post, query, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", &domain.DomainError{Code: constants.NotFound, Message: "post not found"}
		}

		return "", &domain.DomainError{Code: constants.FindError, Message: "error during get post image"}
	}

	return post.ImageName, nil
}
