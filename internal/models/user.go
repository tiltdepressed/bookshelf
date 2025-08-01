package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string `json:"username" gorm:"unique;not null"`
	PasswordHash   string `json:"password_hash" gorm:"not null"` // Хранить хеш!
	Role           string `json:"role" gorm:"default:user"`
	FavouriteBooks []Book `gorm:"many2many:favourite_books;"`
}
