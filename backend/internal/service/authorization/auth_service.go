package authorization

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/domain"
	"myapp/internal/repository"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

type AuthServiceInterface interface {
	Auth(userDTO dto.UserDTO, ctx context.Context) (*AuthResult, *domain.DomainError)
	Refresh(oldRefreshToken string, ctx context.Context) (*AuthResult, *domain.DomainError)
}

type AuthService struct {
	userRepository         repository.UserRepositoryInterface
	refreshTokenRepository repository.RefreshTokenRepositoryInterface
	jwtService             JWTServiceInterface
	db                     *gorm.DB
}

func NewAuthService(userRepo repository.UserRepositoryInterface, refreshTokenRepo repository.RefreshTokenRepositoryInterface, jwtService JWTServiceInterface, db *gorm.DB) AuthServiceInterface {
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

	if domainErr := s.refreshTokenRepository.Set(strconv.Itoa(int(user.ID)), authResult.RefreshToken, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return nil, domainErr
	}

	return authResult, nil
}

func (s AuthService) Refresh(oldRefreshToken string, ctx context.Context) (*AuthResult, *domain.DomainError) {
	claims, domainErr := s.jwtService.ParseToken(oldRefreshToken)
	if domainErr != nil {
		return nil, domainErr
	}

	dbRefreshToken, domainErr := s.refreshTokenRepository.Get(claims.Subject, ctx)
	if domainErr != nil {
		return nil, domainErr
	}
	if oldRefreshToken != dbRefreshToken {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	authResult, domainErr := s.generateTokenPair(claims.Subject)
	if domainErr != nil {
		return nil, domainErr
	}

	if domainErr := s.refreshTokenRepository.Set(claims.Subject, authResult.RefreshToken, constants.RefreshTokenExpiresAt, ctx); domainErr != nil {
		return nil, domainErr
	}

	return authResult, nil
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
