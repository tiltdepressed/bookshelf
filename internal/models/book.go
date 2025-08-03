// Package models
package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string  `json:"title" gorm:"not null"`
	Author      string  `json:"author" gorm:"not null"`
	Genre       string  `json:"genre" gorm:"not null"`
	Description string  `json:"description" gorm:"not null"`
	Price       float64 `json:"price" gorm:"not null"`
}
