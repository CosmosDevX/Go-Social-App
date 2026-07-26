package dto

import "errors"

type CommentDTO struct {
	CommentID       uint   `json:"comment_id"`
	PostID          uint   `json:"post_id"`
	CommentText     string `json:"comment_text"`
	CreatorUsername string `json:"creator_username"`
	CreatorID       uint   `json:"creator_id"`
}

func (dto CommentDTO) Validate() error {
	if len(dto.CommentText) > 250 || len(dto.CommentText) < 1 {
		return errors.New("comment text cannot be bigger than 250 and lower than 1")
	}

	return nil
}
