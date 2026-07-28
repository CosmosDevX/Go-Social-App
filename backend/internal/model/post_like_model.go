package model

type PostLike struct {
	ID          uint
	LikedUserID uint `gorm:"index"`
	PostID      uint `gorm:"index"`
	Post        Post `gorm:"foreignKey: PostID; references: ID"`
}
