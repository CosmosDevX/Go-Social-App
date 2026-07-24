// Package model
package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type: VARCHAR(60); unique; not null"`
	Password string `gorm:"type: VARCHAR(200); unique; not null"`
}
