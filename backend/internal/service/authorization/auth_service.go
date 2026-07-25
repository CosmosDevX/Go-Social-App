package authorization

import (
	"context"
	"myapp/internal/constants"
	"myapp/internal/delivery"
	"myapp/internal/delivery/http/dto"
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
	Auth(userDTO dto.UserDTO, ctx context.Context) (*AuthResult, *delivery.APIError)
	Refresh(oldRefreshToken string, ctx context.Context) (*AuthResult, *delivery.APIError)
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

func (s AuthService) Auth(userDTO dto.UserDTO, ctx context.Context) (*AuthResult, *delivery.APIError) {
	user, err := s.userRepository.GetUserByName(userDTO.Username, s.db.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password)); err != nil {
		return nil, &delivery.APIError{Code: constants.InvalidPassword, Message: "passwords do not match"}
	}

	authResult, apiErr := s.generateTokenPair(strconv.Itoa(int(user.ID)))
	if apiErr != nil {
		return nil, apiErr
	}

	if apiErr := s.refreshTokenRepository.Set(strconv.Itoa(int(user.ID)), authResult.RefreshToken, constants.RefreshTokenExpiresAt, ctx); apiErr != nil {
		return nil, apiErr
	}

	return authResult, nil
}

func (s AuthService) Refresh(oldRefreshToken string, ctx context.Context) (*AuthResult, *delivery.APIError) {
	claims, apiErr := s.jwtService.ParseToken(oldRefreshToken)
	if apiErr != nil {
		return nil, apiErr
	}

	dbRefreshToken, apiErr := s.refreshTokenRepository.Get(claims.Subject, ctx)
	if apiErr != nil {
		return nil, apiErr
	}
	if oldRefreshToken != dbRefreshToken {
		return nil, &delivery.APIError{Code: constants.InvalidTokenError, Message: "invalid refresh token"}
	}

	authResult, apiErr := s.generateTokenPair(claims.Subject)
	if apiErr != nil {
		return nil, apiErr
	}

	if apiErr := s.refreshTokenRepository.Set(claims.Subject, authResult.RefreshToken, constants.RefreshTokenExpiresAt, ctx); apiErr != nil {
		return nil, apiErr
	}

	return authResult, nil
}

func (s AuthService) generateTokenPair(userID string) (*AuthResult, *delivery.APIError) {
	accessToken, jwtError := s.jwtService.GenerateToken(userID, constants.AccessTokenExpiresAt)
	if jwtError != nil {
		return nil, &delivery.APIError{Code: constants.AccessTokenError, Message: "error during the access token generating"}
	}

	refreshToken, jwtError := s.jwtService.GenerateToken(userID, constants.RefreshTokenExpiresAt)
	if jwtError != nil {
		return nil, &delivery.APIError{Code: constants.RefreshTokenError, Message: "error during the refresh token generating"}
	}

	return &AuthResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
