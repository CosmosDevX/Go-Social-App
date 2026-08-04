package authorization

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

type UserGetter interface {
	GetUserByName(ctx context.Context, username string) (*domain.User, *domain.DomainError)
}

type RefreshTokenRepository interface {
	Set(userID, refreshToken, prefix string, ttl time.Duration, ctx context.Context) *domain.DomainError
	Delete(userID, prefix string, ctx context.Context) *domain.DomainError
	Get(userID, prefix string, ctx context.Context) (string, *domain.DomainError)
}

type AuthService struct {
	userGetter             UserGetter
	refreshTokenRepository RefreshTokenRepository
	jwtService             JWTServiceInterface
}

func NewAuthService(userGetter UserGetter, refreshTokenRepo RefreshTokenRepository, jwtService JWTServiceInterface) AuthService {
	return AuthService{
		userGetter:             userGetter,
		refreshTokenRepository: refreshTokenRepo,
		jwtService:             jwtService,
	}
}

func (s AuthService) Auth(ctx context.Context, userDTO dto.UserDTO) (*AuthResult, *domain.DomainError) {
	user, err := s.userGetter.GetUserByName(ctx, userDTO.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		return nil, &domain.DomainError{Code: constants.InvalidPassword, Message: "passwords do not match"}
	}

	authResult, domainErr := s.generateTokenPair(user.ID, user.Username)
	if domainErr != nil {
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(strconv.Itoa(int(user.ID)), authResult.RefreshToken, constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return nil, domainErr
	}

	return authResult, nil
}

func (s AuthService) Refresh(ctx context.Context, oldRefreshToken string) (*AuthResult, *domain.DomainError) {
	claims, domainErr := s.jwtService.ParseToken(oldRefreshToken)
	if domainErr != nil {
		return nil, domainErr
	}

	dbRefreshToken, domainErr := s.refreshTokenRepository.Get(strconv.Itoa(claims.UserID), constants.TokenWhiteListPrefix, ctx)
	if domainErr != nil {
		return nil, domainErr
	}
	if oldRefreshToken != dbRefreshToken {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	blacklistRefreshToken, domainErr := s.refreshTokenRepository.Get(strconv.Itoa(claims.UserID), constants.TokenBlackListPrefix, ctx)
	if domainErr != nil {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "error during get refresh token in blacklist"}
	}
	if blacklistRefreshToken == oldRefreshToken {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	authResult, domainErr := s.generateTokenPair(claims.UserID, claims.Username)
	if domainErr != nil {
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(strconv.Itoa(claims.UserID), authResult.RefreshToken, constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return nil, domainErr
	}

	return authResult, nil
}

func (s AuthService) Logout(ctx context.Context, userID int, refreshToken string) *domain.DomainError {
	if _, domainErr := s.refreshTokenRepository.Get(strconv.Itoa(userID), constants.TokenWhiteListPrefix, ctx); domainErr != nil {
		return domainErr
	}

	if domainErr := s.refreshTokenRepository.Delete(strconv.Itoa(userID), constants.TokenWhiteListPrefix, ctx); domainErr != nil {
		return domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(strconv.Itoa(userID), refreshToken, constants.TokenBlackListPrefix, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return domainErr
	}

	return nil
}

func (s AuthService) generateTokenPair(userID int, username string) (*AuthResult, *domain.DomainError) {
	accessToken, jwtError := s.jwtService.GenerateToken(userID, username, constants.AccessTokenExpiresAt)
	if jwtError != nil {
		return nil, &domain.DomainError{Code: constants.AccessTokenError, Message: "error during the access token generating"}
	}

	refreshToken, jwtError := s.jwtService.GenerateToken(userID, username, constants.RefreshTokenExpiresAt)
	if jwtError != nil {
		return nil, &domain.DomainError{Code: constants.RefreshTokenError, Message: "error during the refresh token generating"}
	}

	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
