package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/tests/helpers"
)

func TestCommentService_CreateComment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		dto         dto.CommentDTO
		creatorID   int
		postID      int
		setup       func(repo *helpers.MockCommentRepository)
		wantID      int
		wantErrCode string
	}{
		{
			name: "success: sets creator and post ids then creates",
			dto: dto.CommentDTO{
				CommentText: "hello world",
			},
			creatorID: 5,
			postID:    10,
			setup: func(repo *helpers.MockCommentRepository) {
				repo.On("Create", ctx, mock.MatchedBy(func(d dto.CommentDTO) bool {
					return d.CommentText == "hello world" && d.CreatorID == 5 && d.PostID == 10
				})).Return(99, nil)
			},
			wantID: 99,
		},
		{
			name: "error: repository fails",
			dto: dto.CommentDTO{
				CommentText: "hello",
			},
			creatorID: 1,
			postID:    2,
			setup: func(repo *helpers.MockCommentRepository) {
				repo.On("Create", ctx, mock.AnythingOfType("dto.CommentDTO")).
					Return(0, &domain.DomainError{Code: constants.CreateError, Message: "db error"})
			},
			wantErrCode: constants.CreateError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(helpers.MockCommentRepository)
			tt.setup(repo)

			svc := service.NewCommentService(repo, nil)
			id, domainErr := svc.CreateComment(ctx, tt.dto, tt.creatorID, tt.postID)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Equal(t, 0, id)
			} else {
				require.Nil(t, domainErr)
				assert.Equal(t, tt.wantID, id)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestCommentService_DeleteComment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		commentID   int
		userID      int
		setup       func(repo *helpers.MockCommentRepository)
		wantErrCode string
	}{
		{
			name:      "success",
			commentID: 3,
			userID:    7,
			setup: func(repo *helpers.MockCommentRepository) {
				repo.On("Delete", ctx, 3, 7).Return(nil)
			},
		},
		{
			name:      "error: not found / not owner",
			commentID: 3,
			userID:    7,
			setup: func(repo *helpers.MockCommentRepository) {
				repo.On("Delete", ctx, 3, 7).Return(&domain.DomainError{
					Code:    constants.NotFound,
					Message: "comment not deleted",
				})
			},
			wantErrCode: constants.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(helpers.MockCommentRepository)
			tt.setup(repo)

			svc := service.NewCommentService(repo, nil)
			domainErr := svc.DeleteComment(ctx, tt.commentID, tt.userID)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
			} else {
				assert.Nil(t, domainErr)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestCommentService_GetAllCommentsByPostID(t *testing.T) {
	ctx := context.Background()

	domainComments := []domain.Comment{
		{ID: 1, CommentText: "first", PostID: 10, CreatorID: 5, CreatorUsername: "alice"},
		{ID: 2, CommentText: "second", PostID: 10, CreatorID: 6, CreatorUsername: "bob"},
	}

	tests := []struct {
		name        string
		postID      int
		setup       func(repo *helpers.MockCommentRepository)
		wantLen     int
		wantErrCode string
	}{
		{
			name:   "success: maps domain comments to DTOs",
			postID: 10,
			setup: func(repo *helpers.MockCommentRepository) {
				repo.On("GetAllByPostID", ctx, 10).Return(domainComments, nil)
			},
			wantLen: 2,
		},
		{
			name:   "error: not found",
			postID: 999,
			setup: func(repo *helpers.MockCommentRepository) {
				repo.On("GetAllByPostID", ctx, 999).Return(nil, &domain.DomainError{
					Code:    constants.NotFound,
					Message: "comments not found",
				})
			},
			wantErrCode: constants.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(helpers.MockCommentRepository)
			tt.setup(repo)

			svc := service.NewCommentService(repo, nil)
			dtos, domainErr := svc.GetAllCommentsByPostID(ctx, tt.postID)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Nil(t, dtos)
			} else {
				require.Nil(t, domainErr)
				require.Len(t, dtos, tt.wantLen)
				assert.Equal(t, "first", dtos[0].CommentText)
				assert.Equal(t, "alice", dtos[0].CreatorUsername)
				assert.Equal(t, "second", dtos[1].CommentText)
			}
			repo.AssertExpectations(t)
		})
	}
}
