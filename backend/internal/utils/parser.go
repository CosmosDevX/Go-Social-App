package utils

import (
	"errors"
	"strconv"
)

func ParseUserID(ctxValue any) (int, error) {
	stringUserID, ok := ctxValue.(string)
	if !ok {
		return 0, errors.New("error during parsing userID to string")
	}

	userID, err := strconv.Atoi(stringUserID)
	if err != nil {
		return 0, errors.New("error during parsing userID to iint")
	}

	return userID, nil
}
