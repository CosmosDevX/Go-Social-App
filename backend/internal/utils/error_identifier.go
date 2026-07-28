package utils

import (
	"myapp/internal/constants"
	"net/http"
)

func IdentifyAPIError(code string) int {
	switch code {
	case constants.NotFound:
		return http.StatusNotFound
	case constants.RequestTimeout:
		return http.StatusGatewayTimeout
	case constants.InvalidPassword:
		return http.StatusUnauthorized
	case constants.TokenSignError:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
