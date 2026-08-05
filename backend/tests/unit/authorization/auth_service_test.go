package authorization_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/service/authorization"
	"myapp/tests/helpers"
)

func TestAuthService_Auth(t *testing.T) {
	ctx := context.Background()
	hashed, err := bcrypt.GenerateFromPassword([]byte("password1234"), 10)
	require.NoError(t, err)

	validUser := &domain.User{
		ID:       42,
		Username: "alice",
		Password: string(hashed),
	}

	tests := []struct {
		name          string
		userDTO       dto.UserDTO
		setup         func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService)
		wantAccess    string
		wantRefresh   string
		wantErrCode   string
		wantErrMsgSub string
	}{
		{
			name: "success: valid credentials returns token pair and stores refresh in whitelist",
			userDTO: dto.UserDTO{
				Username: "alice",
				Password: "password1234",
			},
			setup: func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				ug.On("GetUserByName", ctx, "alice").Return(validUser, nil)
				jwt.On("GenerateToken", 42, "alice", constants.AccessTokenExpiresAt).Return("access-token", nil)
				jwt.On("GenerateToken", 42, "alice", constants.RefreshTokenExpiresAt).Return("refresh-token", nil)
				rt.On("Set", "42", "refresh-token", constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx).Return(nil)
			},
			wantAccess:  "access-token",
			wantRefresh: "refresh-token",
		},
		{
			name: "error: user not found",
			userDTO: dto.UserDTO{
				Username: "unknown",
				Password: "password1234",
			},
			setup: func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				ug.On("GetUserByName", ctx, "unknown").Return(nil, &domain.DomainError{
					Code:    constants.NotFound,
					Message: "user not found",
				})
			},
			wantErrCode: constants.NotFound,
		},
		{
			name: "error: wrong password",
			userDTO: dto.UserDTO{
				Username: "alice",
				Password: "wrong-password",
			},
			setup: func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				ug.On("GetUserByName", ctx, "alice").Return(validUser, nil)
			},
			wantErrCode: constants.InvalidPassword,
		},
		{
			name: "error: access token generation fails",
			userDTO: dto.UserDTO{
				Username: "alice",
				Password: "password1234",
			},
			setup: func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				ug.On("GetUserByName", ctx, "alice").Return(validUser, nil)
				jwt.On("GenerateToken", 42, "alice", constants.AccessTokenExpiresAt).Return("", &domain.DomainError{
					Code:    constants.TokenSignError,
					Message: "sign error",
				})
			},
			wantErrCode: constants.AccessTokenError,
		},
		{
			name: "error: refresh token generation fails",
			userDTO: dto.UserDTO{
				Username: "alice",
				Password: "password1234",
			},
			setup: func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				ug.On("GetUserByName", ctx, "alice").Return(validUser, nil)
				jwt.On("GenerateToken", 42, "alice", constants.AccessTokenExpiresAt).Return("access-token", nil)
				jwt.On("GenerateToken", 42, "alice", constants.RefreshTokenExpiresAt).Return("", &domain.DomainError{
					Code:    constants.TokenSignError,
					Message: "sign error",
				})
			},
			wantErrCode: constants.RefreshTokenError,
		},
		{
			name: "error: failed to store refresh token in whitelist",
			userDTO: dto.UserDTO{
				Username: "alice",
				Password: "password1234",
			},
			setup: func(ug *helpers.MockUserGetter, rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				ug.On("GetUserByName", ctx, "alice").Return(validUser, nil)
				jwt.On("GenerateToken", 42, "alice", constants.AccessTokenExpiresAt).Return("access-token", nil)
				jwt.On("GenerateToken", 42, "alice", constants.RefreshTokenExpiresAt).Return("refresh-token", nil)
				rt.On("Set", "42", "refresh-token", constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx).
					Return(&domain.DomainError{Code: constants.SaveError, Message: "redis error"})
			},
			wantErrCode: constants.SaveError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ug := new(helpers.MockUserGetter)
			rt := new(helpers.MockRefreshTokenRepository)
			jwt := new(helpers.MockJWTService)
			tt.setup(ug, rt, jwt)

			svc := authorization.NewAuthService(ug, rt, jwt)
			result, domainErr := svc.Auth(ctx, tt.userDTO)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Nil(t, result)
			} else {
				require.Nil(t, domainErr)
				require.NotNil(t, result)
				assert.Equal(t, tt.wantAccess, result.AccessToken)
				assert.Equal(t, tt.wantRefresh, result.RefreshToken)
			}

			ug.AssertExpectations(t)
			rt.AssertExpectations(t)
			jwt.AssertExpectations(t)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	ctx := context.Background()
	claims := &authorization.UserClaims{
		UserID:   42,
		Username: "alice",
	}

	tests := []struct {
		name        string
		oldToken    string
		setup       func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService)
		wantAccess  string
		wantRefresh string
		wantErrCode string
	}{
		{
			name:     "success: valid refresh token rotates tokens",
			oldToken: "old-refresh",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "old-refresh").Return(claims, nil)
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("old-refresh", nil)
				rt.On("Get", "42", constants.TokenBlackListPrefix, ctx).Return("", nil)
				jwt.On("GenerateToken", 42, "alice", constants.AccessTokenExpiresAt).Return("new-access", nil)
				jwt.On("GenerateToken", 42, "alice", constants.RefreshTokenExpiresAt).Return("new-refresh", nil)
				rt.On("Set", "42", "new-refresh", constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx).Return(nil)
			},
			wantAccess:  "new-access",
			wantRefresh: "new-refresh",
		},
		{
			name:     "error: parse token fails",
			oldToken: "bad-token",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "bad-token").Return(nil, &domain.DomainError{
					Code:    constants.InvalidTokenError,
					Message: "invalid token",
				})
			},
			wantErrCode: constants.InvalidTokenError,
		},
		{
			name:     "error: whitelist token does not match (stolen or rotated)",
			oldToken: "old-refresh",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "old-refresh").Return(claims, nil)
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("different-token", nil)
			},
			wantErrCode: constants.InvalidTokenError,
		},
		{
			name:     "error: token is blacklisted",
			oldToken: "old-refresh",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "old-refresh").Return(claims, nil)
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("old-refresh", nil)
				rt.On("Get", "42", constants.TokenBlackListPrefix, ctx).Return("old-refresh", nil)
			},
			wantErrCode: constants.InvalidTokenError,
		},
		{
			name:     "error: get from whitelist fails",
			oldToken: "old-refresh",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "old-refresh").Return(claims, nil)
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("", &domain.DomainError{
					Code:    constants.FindError,
					Message: "redis error",
				})
			},
			wantErrCode: constants.FindError,
		},
		{
			name:     "error: get from blacklist returns error (treated as invalid)",
			oldToken: "old-refresh",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "old-refresh").Return(claims, nil)
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("old-refresh", nil)
				rt.On("Get", "42", constants.TokenBlackListPrefix, ctx).Return("", &domain.DomainError{
					Code:    constants.FindError,
					Message: "redis error",
				})
			},
			wantErrCode: constants.InvalidTokenError,
		},
		{
			name:     "error: failed to store new refresh token",
			oldToken: "old-refresh",
			setup: func(rt *helpers.MockRefreshTokenRepository, jwt *helpers.MockJWTService) {
				jwt.On("ParseToken", "old-refresh").Return(claims, nil)
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("old-refresh", nil)
				rt.On("Get", "42", constants.TokenBlackListPrefix, ctx).Return("", nil)
				jwt.On("GenerateToken", 42, "alice", constants.AccessTokenExpiresAt).Return("new-access", nil)
				jwt.On("GenerateToken", 42, "alice", constants.RefreshTokenExpiresAt).Return("new-refresh", nil)
				rt.On("Set", "42", "new-refresh", constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx).
					Return(&domain.DomainError{Code: constants.SaveError, Message: "redis down"})
			},
			wantErrCode: constants.SaveError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := new(helpers.MockRefreshTokenRepository)
			jwt := new(helpers.MockJWTService)
			tt.setup(rt, jwt)

			svc := authorization.NewAuthService(nil, rt, jwt)
			result, domainErr := svc.Refresh(ctx, tt.oldToken)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
				assert.Nil(t, result)
			} else {
				require.Nil(t, domainErr)
				require.NotNil(t, result)
				assert.Equal(t, tt.wantAccess, result.AccessToken)
				assert.Equal(t, tt.wantRefresh, result.RefreshToken)
			}

			rt.AssertExpectations(t)
			jwt.AssertExpectations(t)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      int
		token       string
		setup       func(rt *helpers.MockRefreshTokenRepository)
		wantErrCode string
	}{
		{
			name:   "success: removes from whitelist and adds to blacklist",
			userID: 42,
			token:  "refresh-token",
			setup: func(rt *helpers.MockRefreshTokenRepository) {
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("refresh-token", nil)
				rt.On("Delete", "42", constants.TokenWhiteListPrefix, ctx).Return(nil)
				rt.On("Set", "42", "refresh-token", constants.TokenBlackListPrefix, constants.RefreshTokenExpiresAt, ctx).Return(nil)
			},
		},
		{
			name:   "error: get from whitelist fails",
			userID: 42,
			token:  "refresh-token",
			setup: func(rt *helpers.MockRefreshTokenRepository) {
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("", &domain.DomainError{
					Code:    constants.FindError,
					Message: "not found",
				})
			},
			wantErrCode: constants.FindError,
		},
		{
			name:   "error: delete from whitelist fails",
			userID: 42,
			token:  "refresh-token",
			setup: func(rt *helpers.MockRefreshTokenRepository) {
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("refresh-token", nil)
				rt.On("Delete", "42", constants.TokenWhiteListPrefix, ctx).Return(&domain.DomainError{
					Code:    constants.DeleteError,
					Message: "redis error",
				})
			},
			wantErrCode: constants.DeleteError,
		},
		{
			name:   "error: set to blacklist fails",
			userID: 42,
			token:  "refresh-token",
			setup: func(rt *helpers.MockRefreshTokenRepository) {
				rt.On("Get", "42", constants.TokenWhiteListPrefix, ctx).Return("refresh-token", nil)
				rt.On("Delete", "42", constants.TokenWhiteListPrefix, ctx).Return(nil)
				rt.On("Set", "42", "refresh-token", constants.TokenBlackListPrefix, constants.RefreshTokenExpiresAt, ctx).
					Return(&domain.DomainError{Code: constants.SaveError, Message: "redis error"})
			},
			wantErrCode: constants.SaveError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := new(helpers.MockRefreshTokenRepository)
			tt.setup(rt)

			svc := authorization.NewAuthService(nil, rt, nil)
			domainErr := svc.Logout(ctx, tt.userID, tt.token)

			if tt.wantErrCode != "" {
				require.NotNil(t, domainErr)
				assert.Equal(t, tt.wantErrCode, domainErr.Code)
			} else {
				assert.Nil(t, domainErr)
			}

			rt.AssertExpectations(t)
		})
	}
}

var _ = time.Second
var _ = mock.Anything
