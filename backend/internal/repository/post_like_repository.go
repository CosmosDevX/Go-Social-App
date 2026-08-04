package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/domain"
)

type PostLikeRepository struct {
	db DBTX
}

func NewPostLikeRepository(db DBTX) PostLikeRepository {
	return PostLikeRepository{
		db: db,
	}
}

func (r PostLikeRepository) GetLikedPostsID(ctx context.Context, userID int) ([]int, *domain.DomainError) {
	query := `SELECT post_ID FROM post_likes WHERE liked_user_id = $1`
	var postLikes []domain.PostLike
	err := r.db.SelectContext(ctx, &postLikes, query, userID)
	if err != nil {
		return nil, &domain.DomainError{
			Code:    constants.FindError,
			Message: "error during finding liked post ids",
		}
	}

	postIDs := make([]int, len(postLikes))
	for i := range postLikes {
		postIDs[i] = postLikes[i].PostID
	}

	return postIDs, nil
}

func (r PostLikeRepository) CreateLike(ctx context.Context, likedUserID, postID int) (int, *domain.DomainError) {
	query := `
		INSERT INTO post_likes(liked_user_id, post_id) VALUES($1, $2) 
		ON CONFLICT (liked_user_id, post_id) DO NOTHING
	`
	sqlResult, err := r.db.ExecContext(ctx, query, likedUserID, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during create post like"}
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return 0, &domain.DomainError{Code: constants.DatabaseError, Message: "error during get affected rows"}
	}

	return int(rowsAffected), nil
}

func (r PostLikeRepository) DeleteLike(ctx context.Context, likedUserID, postID int) (int, *domain.DomainError) {
	query := `DELETE FROM post_likes WHERE liked_user_id = $1 AND post_id = $2`
	sqlResult, err := r.db.ExecContext(ctx, query, likedUserID, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.DeleteError, Message: "error during delete post like"}
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return 0, &domain.DomainError{Code: constants.DatabaseError, Message: "error during get affected rows"}
	}

	return int(rowsAffected), nil
}
