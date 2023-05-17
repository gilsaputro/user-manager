package postgres

import "github.com/jinzhu/gorm"

// User struct to user information
type User struct {
	gorm.Model
	Username string `gorm:"size:50;unique;not null"`
	Password string `gorm:"not null"`
	Fullname string
	Email    string
}
