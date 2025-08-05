// Package models
package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model  `swaggerignore:"true"`
	Title       string  `json:"title" gorm:"not null" example:"The Go Programming Language"`
	Author      string  `json:"author" gorm:"not null" example:"Alan A. A. Donovan"`
	Genre       string  `json:"genre" gorm:"not null" example:"Programming"`
	Description string  `json:"description" gorm:"not null" example:"Definitive guide to Go programming"`
	Price       float64 `json:"price" gorm:"not null" example:"49.99"`
}
