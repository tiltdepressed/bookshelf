// Package repository
package repository

import (
	models "bookshelf/internal/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user models.User) error
	GetAllUsers() ([]models.User, error)
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(id string) (models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) error
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(user models.User) error {
	return r.db.Create(&user).Error
}

func (r *authRepo) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *authRepo) UpdateUser(user models.User) error {
	return r.db.Save(&user).Error
}

func (r *authRepo) DeleteUser(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *authRepo) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *authRepo) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	return user, err
}
