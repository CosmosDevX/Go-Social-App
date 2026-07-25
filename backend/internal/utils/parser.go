package utils

import (
	"errors"
	"strconv"
)

func ParseUserID(ctxValue any) (uint, error) {
	stringUserID, ok := ctxValue.(string)
	if !ok {
		return 0, errors.New("error during parsing userID to string")
	}

	userID, err := strconv.ParseUint(stringUserID, 10, 64)
	if err != nil {
		return 0, errors.New("error during parsing userID to uint")
	}

	return uint(userID), nil
}
