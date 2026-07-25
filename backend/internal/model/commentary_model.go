package model

type Commentary struct {
	ID             uint
	CommentaryText string `gorm:"type: VARCHAR(250)"`
	PostID         uint   `gorm:"index"`
	Post           Post   `gorm:"foreignKey: PostID; references: ID"`
	CreatorID      uint   `gorm:"index"`
	Creator        User   `gorm:"foreignKey: CreatorID; references: ID"`
}
