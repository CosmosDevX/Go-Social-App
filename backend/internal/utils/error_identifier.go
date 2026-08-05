package utils

import (
	"myapp/internal/constants"
	"net/http"
)

func MapError(code string) int {
	switch code {
	case constants.NotFound:
		return http.StatusNotFound
	case constants.RequestTimeout:
		return http.StatusGatewayTimeout
	case constants.InvalidPassword:
		return http.StatusUnauthorized
	case constants.TokenSignError, constants.InvalidTokenError, constants.AuthError:
		return http.StatusUnauthorized
	case constants.ValidationError, constants.DeserializingError, constants.FileError, constants.ParseError:
		return http.StatusBadRequest
	case constants.TooManyRequests:
		return http.StatusTooManyRequests
	case constants.UniqueViolation:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
