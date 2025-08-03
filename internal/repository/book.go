package repository

import (
	models "bookshelf/internal/models"

	"gorm.io/gorm"
)

type BookRepository interface {
	CreateBook(book models.Book) error
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id string) (models.Book, error)
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

func (r *bookRepo) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.db.First(&books).Error
	return books, err
}

func (r *bookRepo) GetBookByID(id string) (models.Book, error) {
	var book models.Book
	err := r.db.First(&book, "id = ?", id).Error
	return book, err
}

func (r *bookRepo) UpdateBook(book models.Book) error {
	return r.db.Save(&book).Error
}

func (r *bookRepo) DeleteBook(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Book{}).Error
}
