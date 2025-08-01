// Package db
package db

import (
	"bookshelf/internal/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to bd: %s", err.Error())
	}
	if err = db.AutoMigrate(&models.Book{}, &models.User{}); err != nil {
		log.Fatalf("Could not migrate: %s", err.Error())
	}

	return db, nil
}
