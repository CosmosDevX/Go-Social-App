// Package authorization
package authorization

import (
	"errors"
	"myapp/internal/constants"
	"myapp/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID   int
	Username string
	jwt.RegisteredClaims
}

type JWTServiceInterface interface {
	GenerateToken(userID int, username string, expiresAt time.Duration) (string, *domain.DomainError)
	ParseToken(tokenString string) (*UserClaims, *domain.DomainError)
}

type JWTService struct {
	secretKey []byte
}

func NewJWTService(secretKey string) JWTServiceInterface {
	return JWTService{
		secretKey: []byte(secretKey),
	}
}

func (s JWTService) GenerateToken(userID int, username string, expiresAt time.Duration) (string, *domain.DomainError) {
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "myapp",
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", &domain.DomainError{Code: constants.TokenSignError, Message: "error during the signing token"}
	}

	return tokenString, nil
}

func (s JWTService) ParseToken(tokenString string) (*UserClaims, *domain.DomainError) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect encrypt method")
		}

		return s.secretKey, nil
	})

	if err != nil {
		return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "invalid token"}
	}

	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return userClaims, nil
	}

	return nil, &domain.DomainError{Code: constants.InvalidTokenError, Message: "token is invalid"}
}
