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
	case constants.TokenSignError:
		return http.StatusUnauthorized
	case constants.ValidationError:
		return http.StatusBadRequest
	case constants.DeserializingError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
