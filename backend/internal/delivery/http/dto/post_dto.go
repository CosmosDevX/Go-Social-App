package dto

import "errors"

type PostDTO struct {
	PostID          uint   `json:"post_id"`
	PostName        string `json:"post_name"`
	PostDescription string `json:"post_description"`
	CreatorID       uint   `json:"creator_id"`
	CreatorName     string `json:"creator_name"`
	Likes           int    `json:"likes"`
	IsLiked         bool   `json:"is_liked"`
	CommentsCount   int    `json:"comments_count"`
	ImageName       string `json:"image_name"`
}

func (dto PostDTO) Validate() error {
	if len(dto.PostName) > 100 || len(dto.PostName) < 5 {
		return errors.New("post name cannot be bigger than 100 and lower than 5")
	}

	if len(dto.PostDescription) > 900 || len(dto.PostDescription) < 1 {
		return errors.New("post description cannot be bigger than 900 and lower than 1")
	}

	return nil
}
