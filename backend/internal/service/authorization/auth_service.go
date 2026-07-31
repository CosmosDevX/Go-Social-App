package authorization

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

type UserGetter interface {
	GetUserByName(username string, db *gorm.DB) (*domain.User, *domain.DomainError)
}

type RefreshTokenRepository interface {
	Set(userID, refreshToken, prefix string, ttl time.Duration, ctx context.Context) *domain.DomainError
	Delete(userID, prefix string, ctx context.Context) *domain.DomainError
	Get(userID, prefix string, ctx context.Context) (string, *domain.DomainError)
}

type AuthService struct {
	userRepository         UserGetter
	refreshTokenRepository RefreshTokenRepository
	jwtService             JWTServiceInterface
	db                     *gorm.DB
}

func NewAuthService(userRepo UserGetter, refreshTokenRepo RefreshTokenRepository, jwtService JWTServiceInterface, db *gorm.DB) AuthService {
	return AuthService{
		userRepository:         userRepo,
		refreshTokenRepository: refreshTokenRepo,
		jwtService:             jwtService,
		db:                     db,
	}
}

func (s AuthService) Auth(userDTO dto.UserDTO, ctx context.Context) (*AuthResult, *domain.DomainError) {
	user, err := s.userRepository.GetUserByName(userDTO.Username, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		return nil, &domain.DomainError{Code: constants.InvalidPassword, Message: "passwords do not match"}
	}

	authResult, domainErr := s.generateTokenPair(strconv.Itoa(int(user.ID)))
	if domainErr != nil {
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(strconv.Itoa(int(user.ID)), authResult.RefreshToken, constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return nil, domainErr
	}

	return authResult, nil
}

func (s AuthService) Refresh(oldRefreshToken string, ctx context.Context) (*AuthResult, *domain.DomainError) {
	claims, domainErr := s.jwtService.ParseToken(oldRefreshToken)
	if domainErr != nil {
		return nil, domainErr
	}

	dbRefreshToken, domainErr := s.refreshTokenRepository.Get(claims.Subject, constants.TokenWhiteListPrefix, ctx)
	if domainErr != nil {
		return nil, domainErr
	}
	if oldRefreshToken != dbRefreshToken {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	blacklistRefreshToken, domainErr := s.refreshTokenRepository.Get(claims.Subject, constants.TokenBlackListPrefix, ctx)
	if domainErr != nil {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "error during get refresh token in blacklist"}
	}
	if blacklistRefreshToken == oldRefreshToken {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	authResult, domainErr := s.generateTokenPair(claims.Subject)
	if domainErr != nil {
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(claims.Subject, authResult.RefreshToken, constants.TokenWhiteListPrefix, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return nil, domainErr
	}

	return authResult, nil
}

func (s AuthService) Logout(userID, refreshToken string, ctx context.Context) *domain.DomainError {
	if _, domainErr := s.refreshTokenRepository.Get(userID, constants.TokenWhiteListPrefix, ctx); domainErr != nil {
		return domainErr
	}

	if domainErr := s.refreshTokenRepository.Delete(userID, constants.TokenWhiteListPrefix, ctx); domainErr != nil {
		return domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(userID, refreshToken, constants.TokenBlackListPrefix, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return domainErr
	}

	return nil
}

func (s AuthService) generateTokenPair(userID string) (*AuthResult, *domain.DomainError) {
	accessToken, jwtError := s.jwtService.GenerateToken(userID, constants.AccessTokenExpiresAt)
	if jwtError != nil {
		return nil, &domain.DomainError{Code: constants.AccessTokenError, Message: "error during the access token generating"}
	}

	refreshToken, jwtError := s.jwtService.GenerateToken(userID, constants.RefreshTokenExpiresAt)
	if jwtError != nil {
		return nil, &domain.DomainError{Code: constants.RefreshTokenError, Message: "error during the refresh token generating"}
	}

	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
