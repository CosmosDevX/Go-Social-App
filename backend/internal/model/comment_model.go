package model

import "myapp/internal/delivery/http/dto"

type Comment struct {
	ID          uint
	CommentText string `gorm:"type: VARCHAR(250)"`
	PostID      uint   `gorm:"index"`
	Post        Post   `gorm:"foreignKey: PostID; references: ID"`
	CreatorID   uint   `gorm:"index"`
	Creator     User   `gorm:"foreignKey: CreatorID; references: ID"`
}

func (c Comment) ToCommentDTO() dto.CommentDTO {
	return dto.CommentDTO{
		PostID:          c.PostID,
		CommentText:     c.CommentText,
		CreatorUsername: c.Creator.Username,
		CreatorID:       c.CreatorID,
	}
}
