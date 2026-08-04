package domain

import "myapp/internal/delivery/http/dto"

type Comment struct {
	ID              int    `db:"id"`
	CommentText     string `db:"text"`
	PostID          int    `db:"post_id"`
	CreatorID       int    `db:"creator_id"`
	CreatorUsername string `db:"creator_username"`
}

func (c Comment) ToCommentDTO() dto.CommentDTO {
	return dto.CommentDTO{
		CommentID:       c.ID,
		PostID:          c.PostID,
		CommentText:     c.CommentText,
		CreatorUsername: c.CreatorUsername,
		CreatorID:       c.CreatorID,
	}
}
