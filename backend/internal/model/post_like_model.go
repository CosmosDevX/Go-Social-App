package model

import "gorm.io/gorm"

type PostLike struct {
	gorm.Model
	LikedUserID uint `gorm:"index"`
	PostID      uint `gorm:"index"`
	Post        Post `gorm:"foreignKey: PostID; references: ID"`
}
