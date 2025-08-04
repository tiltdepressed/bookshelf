package repository

import (
	models "bookshelf/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBook(book models.Book) error
	GetAllBooks(genre string, page, limit int) ([]models.Book, int64, error)
	GetBookByID(id string) (models.Book, error)
	GetAllGenres() ([]string, error)
	UpdateBook(book models.Book) error
	DeleteBook(id string) error
}

type bookRepo struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepo{db: db}
}

func (r *bookRepo) CreateBook(book models.Book) error {
	return r.db.Create(&book).Error
}

func (r *bookRepo) GetAllBooks(genre string, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	db := r.db.Model(&models.Book{})

	if genre != "" {
		db = db.Where("genre = ?", genre)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := db.Offset(offset).Limit(limit).Find(&books).Error
	return books, total, err
}

func (r *bookRepo) GetBookByID(id string) (models.Book, error) {
	var book models.Book
	err := r.db.First(&book, "id = ?", id).Error
	return book, err
}

func (r *bookRepo) GetAllGenres() ([]string, error) {
	var genres []string
	err := r.db.Model(&models.Book{}).Distinct("genre").Pluck("genre", &genres).Error
	return genres, err
}

func (r *bookRepo) UpdateBook(book models.Book) error {
	return r.db.Save(&book).Error
}

func (r *bookRepo) DeleteBook(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Book{}).Error
}
