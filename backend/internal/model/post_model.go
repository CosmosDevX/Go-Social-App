package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	PostName        string `gorm:"type: VARCHAR(100)"`
	PostDescription string `gorm:"type: VARCHAR(900)"`
	CreatorID       uint   `gorm:"index"`
	Creator         User   `gorm:"foreignKey: CreatorID; references: ID"`
	Likes           int    `gorm:"type: INTEGER"`
}
