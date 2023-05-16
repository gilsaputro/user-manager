package postgres

import "github.com/jinzhu/gorm"

// User struct to user information
type User struct {
	gorm.Model
	Username string
	Password string
}
