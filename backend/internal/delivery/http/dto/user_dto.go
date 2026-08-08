// Package dto
package dto

import (
	"errors"
	"unicode/utf8"
)

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (dto UserDTO) Validate() error {
	usernameLength := utf8.RuneCountInString(dto.Username)
	passwordLength := utf8.RuneCountInString(dto.Password)

	if usernameLength > 60 || usernameLength < 3 {
		return errors.New("username cannot be bigger than 60 and lower than 3")
	}

	if passwordLength > 100 || passwordLength < 10 {
		return errors.New("password cannot be bigger than 100 and lower than 10")
	}

	return nil
}
