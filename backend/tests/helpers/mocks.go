package helpers

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/service/authorization"
)

func domainErrFromArgs(args mock.Arguments, idx int) *domain.DomainError {
	if args.Get(idx) == nil {
		return nil
	}
	return args.Get(idx).(*domain.DomainError)
}

// --- AuthService dependencies ---

type MockUserGetter struct {
	mock.Mock
}

func (m *MockUserGetter) GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError) {
	args := m.Called(ctx, username)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}
	return user, domainErrFromArgs(args, 1)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Set(userID, refreshToken, prefix string, ttl time.Duration, ctx context.Context) *domain.DomainError {
	args := m.Called(userID, refreshToken, prefix, ttl, ctx)
	return domainErrFromArgs(args, 0)
}

func (m *MockRefreshTokenRepository) Delete(userID, prefix string, ctx context.Context) *domain.DomainError {
	args := m.Called(userID, prefix, ctx)
	return domainErrFromArgs(args, 0)
}

func (m *MockRefreshTokenRepository) Get(userID, prefix string, ctx context.Context) (string, *domain.DomainError) {
	args := m.Called(userID, prefix, ctx)
	return args.String(0), domainErrFromArgs(args, 1)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(userID int, username string, expiresAt time.Duration) (string, *domain.DomainError) {
	args := m.Called(userID, username, expiresAt)
	return args.String(0), domainErrFromArgs(args, 1)
}

func (m *MockJWTService) ParseToken(tokenString string) (*authorization.UserClaims, *domain.DomainError) {
	args := m.Called(tokenString)
	var claims *authorization.UserClaims
	if args.Get(0) != nil {
		claims = args.Get(0).(*authorization.UserClaims)
	}
	return claims, domainErrFromArgs(args, 1)
}

// --- UserService ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError) {
	args := m.Called(ctx, username)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}
	return user, domainErrFromArgs(args, 1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, userDTO dto.UserDTO) (int, *domain.DomainError) {
	args := m.Called(ctx, userDTO)
	return args.Int(0), domainErrFromArgs(args, 1)
}

func (m *MockUserRepository) GetUsernameByID(ctx context.Context, userID int) (string, *domain.DomainError) {
	args := m.Called(ctx, userID)
	return args.String(0), domainErrFromArgs(args, 1)
}

// --- PostLikeService dependencies ---

type MockPostLikeRepo struct {
	mock.Mock
}

func (m *MockPostLikeRepo) CreateLike(ctx context.Context, likedUserID, postID int) (int, *domain.DomainError) {
	args := m.Called(ctx, likedUserID, postID)
	return args.Int(0), domainErrFromArgs(args, 1)
}

func (m *MockPostLikeRepo) DeleteLike(ctx context.Context, likedUserID, postID int) (int, *domain.DomainError) {
	args := m.Called(ctx, likedUserID, postID)
	return args.Int(0), domainErrFromArgs(args, 1)
}

type MockLikeUpdater struct {
	mock.Mock
}

func (m *MockLikeUpdater) IncrementLikes(ctx context.Context, postID int) (int, *domain.DomainError) {
	args := m.Called(ctx, postID)
	return args.Int(0), domainErrFromArgs(args, 1)
}

func (m *MockLikeUpdater) DecrementLikes(ctx context.Context, postID int) (int, *domain.DomainError) {
	args := m.Called(ctx, postID)
	return args.Int(0), domainErrFromArgs(args, 1)
}

// --- CommentService ---

type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Delete(ctx context.Context, commentID, userID int) *domain.DomainError {
	args := m.Called(ctx, commentID, userID)
	return domainErrFromArgs(args, 0)
}

func (m *MockCommentRepository) Create(ctx context.Context, commentDTO dto.CommentDTO) (int, *domain.DomainError) {
	args := m.Called(ctx, commentDTO)
	return args.Int(0), domainErrFromArgs(args, 1)
}

func (m *MockCommentRepository) GetAllByPostID(ctx context.Context, postID int) ([]domain.Comment, *domain.DomainError) {
	args := m.Called(ctx, postID)
	var comments []domain.Comment
	if args.Get(0) != nil {
		comments = args.Get(0).([]domain.Comment)
	}
	return comments, domainErrFromArgs(args, 1)
}

type MockUsernamesGetter struct {
	mock.Mock
}

func (m *MockUsernamesGetter) GetUsernameByIDs(ctx context.Context, ids []int) ([]string, *domain.DomainError) {
	args := m.Called(ctx, ids)
	var names []string
	if args.Get(0) != nil {
		names = args.Get(0).([]string)
	}
	return names, domainErrFromArgs(args, 1)
}
