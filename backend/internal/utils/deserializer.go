// Package utils
package utils

import (
	"encoding/json"
	"io"
)

func Deserialize(readCloser io.ReadCloser, data any) error {
	body, err := io.ReadAll(readCloser)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	if err := json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}
