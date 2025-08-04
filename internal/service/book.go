package service

import (
	"bookshelf/internal/models"
	"bookshelf/internal/repository"
	"bookshelf/pkg/cache"
	"fmt"
	"time"
)

type BookRequest struct {
	Title       string  `json:"title" binding:"required"`
	Author      string  `json:"author" binding:"required"`
	Genre       string  `json:"genre" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

type BookBrief struct {
	ID     uint    `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Genre  string  `json:"genre"`
	Price  float64 `json:"price"`
}

type BookService interface {
	CreateBook(book BookRequest) (models.Book, error)
	GetBookByID(id string) (models.Book, error)
	GetAllBooks(genre string, page, limit int) ([]BookBrief, int64, error)
	GetAllGenres() ([]string, error)
	UpdateBook(id string, update BookRequest) (models.Book, error)
	DeleteBook(id string) error
}

type bookService struct {
	repo  repository.BookRepository
	cache cache.RedisCache
}

func NewBookService(repo repository.BookRepository, cache *cache.RedisCache) BookService {
	return &bookService{repo: repo, cache: *cache}
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
	s.cache.InvalidatePattern("books:*")
	return book, nil
}

func (s *bookService) GetBookByID(id string) (models.Book, error) {
	cacheKey := fmt.Sprintf("book:%s", id)
	var book models.Book
	if s.cache.Get(cacheKey, &book) {
		return book, nil
	}

	book, err := s.repo.GetBookByID(id)
	if err != nil {
		return models.Book{}, err
	}
	s.cache.Set(cacheKey, book, 10*time.Minute)
	return book, nil
}

func (s *bookService) GetAllBooks(genre string, page, limit int) ([]BookBrief, int64, error) {
	cacheKey := fmt.Sprintf("books:%s:%d:%d", genre, page, limit)

	var cachedResult struct {
		Books []BookBrief
		Total int64
	}

	if s.cache.Get(cacheKey, cachedResult) {
		return cachedResult.Books, cachedResult.Total, nil
	}

	books, total, err := s.repo.GetAllBooks(genre, page, limit)
	if err != nil {
		return nil, 0, err
	}

	briefs := make([]BookBrief, len(books))
	for i, book := range books {
		briefs[i] = BookBrief{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Genre:  book.Genre,
			Price:  book.Price,
		}
	}

	result := struct {
		Books []BookBrief
		Total int64
	}{briefs, total}

	s.cache.Set(cacheKey, result, 5*time.Minute)

	return briefs, total, nil
}

func (s *bookService) GetAllGenres() ([]string, error) {
	return s.repo.GetAllGenres()
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

	s.cache.InvalidatePattern("books:*")
	s.cache.Delete(fmt.Sprintf("book:%s", id))
	return book, nil
}

func (s *bookService) DeleteBook(id string) error {
	if err := s.repo.DeleteBook(id); err != nil {
		return err
	}
	s.cache.InvalidatePattern("book:*")
	s.cache.Delete(fmt.Sprintf("book:%s", id))
	return nil
}
