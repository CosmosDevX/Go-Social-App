package model

import (
	"myapp/internal/delivery/http/dto"
)

type Post struct {
	ID              uint
	PostName        string `gorm:"type: VARCHAR(100)"`
	PostDescription string `gorm:"type: VARCHAR(900)"`
	CreatorID       uint   `gorm:"index"`
	Creator         User   `gorm:"foreignKey: CreatorID; references: ID"`
	Likes           int    `gorm:"type: INTEGER"`
	ImageName       string `gorm:"type: VARCHAR(300)"`

	PostLikes []PostLike `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
	Comments  []Comment  `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;"`
}

func (p Post) ToPostDTO() dto.PostDTO {
	return dto.PostDTO{
		PostID:          p.ID,
		PostName:        p.PostName,
		PostDescription: p.PostDescription,
		CreatorID:       p.CreatorID,
		Likes:           p.Likes,
		CreatorName:     p.Creator.Username,
		ImageName:       p.ImageName,
	}
}
