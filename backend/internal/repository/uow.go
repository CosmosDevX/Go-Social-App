package repository

import (
	"context"
	"database/sql"
	"log"
	"myapp/internal/constants"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

type Repositories struct {
	PostRepository     PostRepository
	PostLikeRepository PostLikeRepository
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) (any, *domain.DomainError)) (any, *domain.DomainError)
}

type unitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) UnitOfWork {
	return unitOfWork{
		db: db,
	}
}

func (u unitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) (any, *domain.DomainError)) (any, *domain.DomainError) {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, &domain.DomainError{Code: constants.TransactionError, Message: "transaction start failed"}
	}

	repos := Repositories{
		PostRepository:     NewPostRepository(tx),
		PostLikeRepository: NewPostLikeRepository(tx),
	}

	value, domainErr := fn(ctx, repos)
	if domainErr != nil {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
		return nil, &domain.DomainError{Code: constants.TransactionError, Message: "transaction failed"}
	}

	if err := tx.Commit(); err != nil {
		return nil, &domain.DomainError{Code: constants.TransactionError, Message: "transaction commit failed"}
	}

	return value, nil
}
