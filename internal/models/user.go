package models

import "gorm.io/gorm"

type User struct {
	gorm.Model     `swaggerignore:"true"`
	Username       string `json:"username" gorm:"unique;not null" example:"john_doe"`
	PasswordHash   string `json:"-" gorm:"not null"`
	Role           string `json:"role" gorm:"default:user" example:"user"`
	FavouriteBooks []Book `gorm:"many2many:user_favourites;"`
}
