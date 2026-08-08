package dto

import (
	"errors"
	"time"
	"unicode/utf8"
)

type PostDTO struct {
	PostID          int       `json:"post_id"`
	PostName        string    `json:"post_name"`
	PostDescription string    `json:"post_description"`
	CreatorID       int       `json:"creator_id"`
	CreatorName     string    `json:"creator_name"`
	Likes           int       `json:"likes"`
	IsLiked         bool      `json:"is_liked"`
	CommentsCount   int       `json:"comments_count"`
	ImageName       string    `json:"image_name"`
	CreatedAt       time.Time `json:"created_at"`
}

func (dto PostDTO) Validate() error {
	postNameLength := utf8.RuneCountInString(dto.PostName)
	postDescriptionLength := utf8.RuneCountInString(dto.PostDescription)

	if postNameLength > 100 || postNameLength < 5 {
		return errors.New("post name cannot be bigger than 100 and lower than 5")
	}

	if postDescriptionLength > 900 || postDescriptionLength < 1 {
		return errors.New("post description cannot be bigger than 900 and lower than 1")
	}

	return nil
}
