package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) CommentRepository {
	return CommentRepository{
		db: db,
	}
}

func (r CommentRepository) Create(ctx context.Context, commentDTO dto.CommentDTO) (int, *domain.DomainError) {
	query := `INSERT INTO comments (text, post_id, creator_id) VALUES($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query, commentDTO.CommentText, commentDTO.PostID, commentDTO.CreatorID).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during create comment"}
	}

	return id, nil
}

func (r CommentRepository) GetAllByPostID(ctx context.Context, postID int) ([]domain.Comment, *domain.DomainError) {
	query := `
		SELECT c.id, text, post_id, creator_id, u.username AS creator_username FROM comments c
		JOIN users u ON u.id = c.creator_id
		WHERE c.post_id = $1
	`
	var comments []domain.Comment
	err := r.db.SelectContext(ctx, &comments, query, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during get comments by post_id"}
	}

	if len(comments) == 0 {
		return nil, &domain.DomainError{Code: constants.NotFound, Message: "comments not found"}
	}

	return comments, nil
}

func (r CommentRepository) CountCommentsOnPost(ctx context.Context, postID int) (int, *domain.DomainError) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, query, postID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.FindError, Message: "error during count comments on post"}
	}

	return count, nil
}

func (r CommentRepository) Delete(ctx context.Context, commentID, userID int) *domain.DomainError {
	query := `DELETE FROM comments WHERE id = $1 AND creator_id = $2`
	sqlResult, err := r.db.ExecContext(ctx, query, commentID, userID)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return &domain.DomainError{Code: constants.DeleteError, Message: "error during delete comment"}
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return &domain.DomainError{Code: constants.DatabaseError, Message: "error during get affected rows"}
	}

	if rowsAffected == 0 {
		return &domain.DomainError{Code: constants.NotFound, Message: "comment not deleted"}
	}

	return nil
}
