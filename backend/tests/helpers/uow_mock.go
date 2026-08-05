// Package helpers
package helpers

import (
	"context"

	"myapp/internal/domain"
	"myapp/internal/repository"

	"github.com/stretchr/testify/mock"
)

// MockUnitOfWork implements repository.UnitOfWork for unit tests.
type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos repository.Repositories) (any, *domain.DomainError)) (any, *domain.DomainError) {
	args := m.Called(ctx, fn)
	// Optionally execute the fn if the test wants to exercise the closure.
	// By default we just return preconfigured values.
	if args.Get(1) == nil {
		return args.Get(0), nil
	}
	return args.Get(0), args.Get(1).(*domain.DomainError)
}

// ExecutingUnitOfWork runs the provided fn with the given repos (no real transaction).
// Useful when you want to unit-test the logic inside the closure with controlled repos.
type ExecutingUnitOfWork struct {
	Repos repository.Repositories
}

func (u ExecutingUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos repository.Repositories) (any, *domain.DomainError)) (any, *domain.DomainError) {
	return fn(ctx, u.Repos)
}
