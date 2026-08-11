// Package constants
package constants

import "time"

const (
	AccessTokenExpiresAt  = 15 * time.Minute
	RefreshTokenExpiresAt = 24 * time.Hour * 7
	RefreshTokenKey       = "refresh_token"
	RefreshTokenMaxAge    = 3600 * 24 * 7
	UsernameKey           = "username"

	TokenWhiteListPrefix = ":tokensWhiteList"
)
