package utils

import (
	"context"
	"errors"
	"fmt"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
)

func ActivateRateLimiter(ctx context.Context, w http.ResponseWriter, r *http.Request, key string, rateLimiter *redis_rate.Limiter, limit redis_rate.Limit) error {
	ip := GetIP(r)

	res, err := rateLimiter.Allow(ctx, key+ip, limit)
	if err != nil {
		WriteError(w, *domain.NewDomainError(constants.TooManyRequests, "too many requests"))
		return errors.New("too many requests err")
	}

	if res.Allowed == 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", res.RetryAfter.Seconds()))
		WriteError(w, *domain.NewDomainError(constants.TooManyRequests, "too many requests. Try again later"))
		return errors.New("too many requests err")
	}

	return nil
}
