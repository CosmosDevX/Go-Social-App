package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/tests/helpers"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userDTO     dto.UserDTO
		setup       func(repo *helpers.MockUserRepository)
		wantID      int
		wantErrCode string
	}{
		{
			name: "success: hashes password and creates user",
			userDTO: dto.UserDTO{
				Username: "bob",
				Password: "securepassword",
			},
			setup: func(repo *helpers.MockUserRepository) {
				repo.On("CreateUser", ctx, mock.MatchedBy(func(d dto.UserDTO) bool {
					if d.Username != "bob" {
						return false
					}
					// password must be bcrypt hash of original
					return bcrypt.CompareHashAndPassword([]byte(d.Password), []byte("securepassword")) == nil
				})).Return(7, nil)
			},
			wantID: 7,
		},
		{
			name: "error: repository create fails",
			userDTO: dto.UserDTO{
				Username: "bob",
				Password: "securepassword",
			},
			setup: func(repo *helpers.MockUserRepository) {
				repo.On("CreateUser", ctx, mock.AnythingOfType("dto.UserDTO")).
					Return(0, &domain.DomainError{Code: constants.CreateError, Message: "db error"})
			},
			wantErrCode: constants.CreateError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(helpers.MockUserRepository)
			tt.setup(repo)

			svc := service.NewUserService(repo)
			id, domainErr := svc.CreateUser(ctx, tt.userDTO)

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

func TestUserService_CurrentUserProfile(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      int
		setup       func(repo *helpers.MockUserRepository)
		wantName    string
		wantErrCode string
	}{
		{
			name:   "success",
			userID: 5,
			setup: func(repo *helpers.MockUserRepository) {
				repo.On("GetUsernameByID", ctx, 5).Return("alice", nil)
			},
			wantName: "alice",
		},
		{
			name:   "error: not found",
			userID: 999,
			setup: func(repo *helpers.MockUserRepository) {
				repo.On("GetUsernameByID", ctx, 999).Return("", &domain.DomainError{
					Code:    constants.NotFound,
					Message: "username not found",
				})
			},
			wantErrCode: constants.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(helpers.MockUserRepository)
			tt.setup(repo)

			svc := service.NewUserService(repo)
			name, domainErr := svc.CurrentUserProfile(ctx, tt.userID)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Empty(t, name)
			} else {
				require.Nil(t, domainErr)
				assert.Equal(t, tt.wantName, name)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUsernameByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		pathValue   string
		setup       func(repo *helpers.MockUserRepository)
		wantName    string
		wantErrCode string
	}{
		{
			name:      "success",
			pathValue: "10",
			setup: func(repo *helpers.MockUserRepository) {
				repo.On("GetUsernameByID", ctx, 10).Return("charlie", nil)
			},
			wantName: "charlie",
		},
		{
			name:        "error: invalid path value (not a number)",
			pathValue:   "abc",
			setup:       func(repo *helpers.MockUserRepository) {},
			wantErrCode: constants.ParseError,
		},
		{
			name:      "error: user not found",
			pathValue: "999",
			setup: func(repo *helpers.MockUserRepository) {
				repo.On("GetUsernameByID", ctx, 999).Return("", &domain.DomainError{
					Code:    constants.NotFound,
					Message: "username not found",
				})
			},
			wantErrCode: constants.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(helpers.MockUserRepository)
			tt.setup(repo)

			svc := service.NewUserService(repo)
			name, domainErr := svc.GetUsernameByID(ctx, tt.pathValue)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Empty(t, name)
			} else {
				require.Nil(t, domainErr)
				assert.Equal(t, tt.wantName, name)
			}
			repo.AssertExpectations(t)
		})
	}
}
