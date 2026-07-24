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
	default:
		return http.StatusInternalServerError
	}
}
