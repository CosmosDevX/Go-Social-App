package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"myapp/internal/constants"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"myapp/internal/service"
	"myapp/tests/helpers"
)

func TestPostLikeService_LikePost(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		postID      int
		userID      int
		setup       func(uow *helpers.MockUnitOfWork)
		wantLikes   int
		wantErrCode string
	}{
		{
			name:   "success: UoW returns new likes count",
			postID: 10,
			userID: 5,
			setup: func(uow *helpers.MockUnitOfWork) {
				uow.On("Do", ctx, mock.AnythingOfType("func(context.Context, repository.Repositories) (interface {}, *domain.DomainError)")).
					Return(42, nil)
			},
			wantLikes: 42,
		},
		{
			name:   "success: already liked (rowsAffected=0) returns 0 likes from UoW",
			postID: 10,
			userID: 5,
			setup: func(uow *helpers.MockUnitOfWork) {
				uow.On("Do", ctx, mock.Anything).Return(0, nil)
			},
			wantLikes: 0,
		},
		{
			name:   "error: UoW / transaction fails",
			postID: 10,
			userID: 5,
			setup: func(uow *helpers.MockUnitOfWork) {
				uow.On("Do", ctx, mock.Anything).Return(nil, &domain.DomainError{
					Code:    constants.TransactionError,
					Message: "tx failed",
				})
			},
			wantErrCode: constants.TransactionError,
		},
		{
			name:   "error: create like fails inside UoW",
			postID: 10,
			userID: 5,
			setup: func(uow *helpers.MockUnitOfWork) {
				uow.On("Do", ctx, mock.Anything).Return(nil, &domain.DomainError{
					Code:    constants.CreateError,
					Message: "error during create post like",
				})
			},
			wantErrCode: constants.CreateError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uow := new(helpers.MockUnitOfWork)
			tt.setup(uow)

			svc := service.NewPostLikeService(uow, nil)
			likes, domainErr := svc.LikePost(ctx, tt.postID, tt.userID)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Equal(t, 0, likes)
			} else {
				require.Nil(t, domainErr)
				assert.Equal(t, tt.wantLikes, likes)
			}
			uow.AssertExpectations(t)
		})
	}
}

func TestPostLikeService_DislikePost(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		postID      int
		userID      int
		setup       func(uow *helpers.MockUnitOfWork)
		wantLikes   int
		wantErrCode string
	}{
		{
			name:   "success: UoW returns decremented likes",
			postID: 10,
			userID: 5,
			setup: func(uow *helpers.MockUnitOfWork) {
				uow.On("Do", ctx, mock.Anything).Return(41, nil)
			},
			wantLikes: 41,
		},
		{
			name:   "error: delete like fails",
			postID: 10,
			userID: 5,
			setup: func(uow *helpers.MockUnitOfWork) {
				uow.On("Do", ctx, mock.Anything).Return(nil, &domain.DomainError{
					Code:    constants.DeleteError,
					Message: "error during delete post like",
				})
			},
			wantErrCode: constants.DeleteError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uow := new(helpers.MockUnitOfWork)
			tt.setup(uow)

			svc := service.NewPostLikeService(uow, nil)
			likes, domainErr := svc.DislikePost(ctx, tt.postID, tt.userID)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
			} else {
				require.Nil(t, domainErr)
				assert.Equal(t, tt.wantLikes, likes)
			}
			uow.AssertExpectations(t)
		})
	}
}

// Ensure repository import is used for type in mock.AnythingOfType.
var _ = repository.Repositories{}
