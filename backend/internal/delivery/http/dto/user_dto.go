// Package dto
package dto

import "errors"

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (dto UserDTO) Validate() error {
	if len(dto.Username) > 60 || len(dto.Username) < 3 {
		return errors.New("username cannot be bigger than 60 and lower than 3")
	}

	if len(dto.Password) > 100 || len(dto.Password) < 10 {
		return errors.New("password cannot be bigger than 100 and lower than 10")
	}

	return nil
}
