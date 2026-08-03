package utils

import (
	"errors"
)

func ParseUserID(ctxValue any) (int, error) {
	userID, ok := ctxValue.(int)
	if !ok {
		return 0, errors.New("error during parsing userID")
	}

	return userID, nil
}
