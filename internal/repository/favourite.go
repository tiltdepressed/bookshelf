package repository

import (
	"bookshelf/internal/models"

	"gorm.io/gorm"
)

type FavouriteRepository interface {
	AddFavourite(userID, bookID uint) error
	RemoveFavourite(userID, bookID uint) error
	GetFavourites(userID uint, page, limit int) ([]models.Book, int64, error)
}

type favouriteRepo struct {
	db *gorm.DB
}

func NewFavouriteRepository(db *gorm.DB) FavouriteRepository {
	return &favouriteRepo{db: db}
}

func (r *favouriteRepo) AddFavourite(userID, bookID uint) error {
	return r.db.Exec(`
		INSERT INTO user_favourites (user_id, book_id)
		VALUES (?, ?)
		ON CONFLICT DO NOTHING
	`, userID, bookID).Error
}

func (r *favouriteRepo) RemoveFavourite(userID, bookID uint) error {
	return r.db.Exec(`
		DELETE FROM user_favourites
		WHERE user_id = ? AND book_id = ?
	`, userID, bookID).Error
}

func (r *favouriteRepo) GetFavourites(userID uint, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	offset := (page - 1) * limit

	err := r.db.Model(&models.Book{}).
		Joins("JOIN user_favourites ON books.id = user_favourites.book_id").
		Where("user_favourites.user_id = ?", userID).
		Count(&total).
		Offset(offset).
		Limit(limit).
		Find(&books).Error

	return books, total, err
}
