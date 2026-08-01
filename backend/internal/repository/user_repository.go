// Package repository
package repository

import (
	"context"
	"database/sql"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

func (r UserRepository) GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError) {
	query := `SELECT * FROM users WHERE username = $1`
	var user domain.User
	err := r.DB.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &domain.DomainError{Code: constants.NotFound, Message: "user not found"}
		}

		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during get user by name"}
	}

	return &user, nil
}

func (r UserRepository) CreateUser(ctx context.Context, userDTO dto.UserDTO) (int, *domain.DomainError) {
	query := `INSERT INTO users (username, password) VALUES($1, $2) RETURNING id`
	var id int
	err := r.DB.QueryRowContext(ctx, query, userDTO.Username, userDTO.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}

		return 0, &domain.DomainError{Code: constants.CreateError, Message: "error during create user"}
	}

	return id, nil
}

func (r UserRepository) GetUsernameByID(ctx context.Context, userID int) (string, *domain.DomainError) {
	query := `SELECT username FROM users WHERE id = $1`
	var user domain.User
	err := r.DB.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", &domain.DomainError{Code: constants.NotFound, Message: "username not found"}
		}

		return "", &domain.DomainError{Code: constants.FindError, Message: "error during get username by user id"}
	}

	return user.Username, nil
}

func (r UserRepository) GetUsernameByIDs(ctx context.Context, ids []int) ([]string, *domain.DomainError) {
	if len(ids) == 0 {
		return nil, &domain.DomainError{Code: constants.NotFound, Message: "usernames not found"}
	}

	query, args, err := sqlx.In(`SELECT username FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during get usernames by ids"}
	}

	query = r.DB.Rebind(query)

	var users []domain.User
	err = r.DB.SelectContext(ctx, &users, query, args...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		return nil, &domain.DomainError{Code: constants.FindError, Message: "error during get usernames by ids"}
	}

	if len(users) == 0 {
		return nil, &domain.DomainError{Code: constants.NotFound, Message: "usernames not found"}
	}

	usernames := make([]string, len(users))
	for i := range users {
		usernames[i] = users[i].Username
	}

	return usernames, nil
}
