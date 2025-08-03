package service

import (
	"bookshelf/internal/models"
	"bookshelf/internal/repository"
)

type BookRequest struct {
	Title       string  `json:"title" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	Genre       string  `json:"genre" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

type BookService interface {
	CreateBook(book BookRequest) (models.Book, error)
	GetBookByID(id string) (models.Book, error)
	GetAllBooks() ([]models.Book, error)
	UpdateBook(id string, update BookRequest) (models.Book, error)
	DeleteBook(id string) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(req BookRequest) (models.Book, error) {
	book := models.Book{
		Title:       req.Title,
		Author:      req.Author,
		Genre:       req.Genre,
		Description: req.Description,
		Price:       req.Price,
	}

	err := s.repo.CreateBook(book)
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}

func (s *bookService) GetBookByID(id string) (models.Book, error) {
	book, err := s.repo.GetBookByID(id)
	return book, err
}

func (s *bookService) GetAllBooks() ([]models.Book, error) {
	books, err := s.repo.GetAllBooks()
	if err != nil {
		return []models.Book{}, err
	}
	return books, nil
}

func (s *bookService) UpdateBook(id string, update BookRequest) (models.Book, error) {
	book, err := s.repo.GetBookByID(id)
	if err != nil {
		return models.Book{}, err
	}
	book.Title = update.Title
	book.Author = update.Author
	book.Genre = update.Genre
	book.Description = update.Description
	book.Price = update.Price

	if err := s.repo.UpdateBook(book); err != nil {
		return models.Book{}, err
	}
	return book, nil
}

func (s *bookService) DeleteBook(id string) error {
	return s.repo.DeleteBook(id)
}
