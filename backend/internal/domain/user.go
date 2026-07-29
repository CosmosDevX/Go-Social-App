// Package domain
package domain

type User struct {
	ID       uint
	Username string `gorm:"type: VARCHAR(60); unique; not null"`
	Password string `gorm:"type: VARCHAR(200); not null"`
}
