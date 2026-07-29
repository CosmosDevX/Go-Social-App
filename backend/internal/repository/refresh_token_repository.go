// Package repository
package repository

import (
	"context"
	"errors"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepositoryInterface interface {
	Set(username, refreshToken string, ttl time.Duration, ctx context.Context) *domain.DomainError
	Delete(username string, ctx context.Context) *domain.DomainError
	Get(username string, ctx context.Context) (string, *domain.DomainError)
}

type RefreshTokenRepository struct {
	redisClient *redis.Client
}

func NewRefreshTokenRepository(redisClient *redis.Client) RefreshTokenRepositoryInterface {
	return RefreshTokenRepository{
		redisClient: redisClient,
	}
}

func (r RefreshTokenRepository) Set(username, refreshToken string, ttl time.Duration, ctx context.Context) *domain.DomainError {
	err := r.redisClient.Set(ctx, username+":tokensWhiteList", refreshToken, ttl).Err() //TODO: Make constant for prefix
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		return &domain.DomainError{Code: constants.SaveError, Message: "error during saving refresh token"}
	}

	return nil
}

func (r RefreshTokenRepository) Delete(username string, ctx context.Context) *domain.DomainError {
	err := r.redisClient.Del(ctx, username+":tokensWhiteList").Err() //TODO: Make constant for prefix
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		return &domain.DomainError{Code: constants.DeleteError, Message: "error during deleting refresh token"}
	}

	return nil
}

func (r RefreshTokenRepository) Get(username string, ctx context.Context) (string, *domain.DomainError) {
	cmd := r.redisClient.Get(ctx, username+":tokensWhiteList") //TODO: Make constant for prefix
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), context.DeadlineExceeded) {
			return "", &domain.DomainError{Code: constants.RequestTimeout, Message: "request timeout"}
		}
		return "", &domain.DomainError{Code: constants.FindError, Message: "no matches found for this username"}
	}

	return cmd.Val(), nil
}
