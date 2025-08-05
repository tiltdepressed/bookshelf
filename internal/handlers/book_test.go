package handlers

import (
	"bookshelf/internal/models"
	"bookshelf/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(book service.BookRequest) (models.Book, error) {
	args := m.Called(book)
	return args.Get(0).(models.Book), args.Error(1)
}

func (m *MockBookService) GetBookByID(id string) (models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(models.Book), args.Error(1)
}

func (m *MockBookService) GetAllBooks(genre string, page, limit int) ([]service.BookBrief, int64, error) {
	args := m.Called(genre, page, limit)
	return args.Get(0).([]service.BookBrief), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookService) GetAllGenres() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockBookService) UpdateBook(id string, update service.BookRequest) (models.Book, error) {
	args := m.Called(id, update)
	return args.Get(0).(models.Book), args.Error(1)
}

func (m *MockBookService) DeleteBook(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBookHandler_CreateBookHandler_Success(t *testing.T) {
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	// Настройка мока
	bookReq := service.BookRequest{
		Title:       "Test Book",
		Author:      "Author",
		Genre:       "Fiction",
		Description: "Description",
		Price:       19.99,
	}
	createdBook := models.Book{
		Model:       gorm.Model{ID: 1},
		Title:       "Test Book",
		Author:      "Author",
		Genre:       "Fiction",
		Description: "Description",
		Price:       19.99,
	}
	mockService.On("CreateBook", bookReq).Return(createdBook, nil)

	// Создание запроса
	bodyBytes, _ := json.Marshal(bookReq)
	req, _ := http.NewRequest("POST", "/books", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.CreateBookHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := `{
		"id":1,
		"title":"Test Book",
		"author":"Author",
		"genre":"Fiction",
		"description":"Description",
		"price":19.99
	}`
	assert.JSONEq(t, expected, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestBookHandler_GetBookByIDHandler_Success(t *testing.T) {
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	// Настройка мока
	book := models.Book{
		Model:       gorm.Model{ID: 1},
		Title:       "Test Book",
		Author:      "Author",
		Genre:       "Fiction",
		Description: "Description",
		Price:       19.99,
	}
	mockService.On("GetBookByID", "1").Return(book, nil)

	// Создание запроса
	req, _ := http.NewRequest("GET", "/books/1", nil)

	// Добавление параметра в роут
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.GetBookByIDHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{
		"id":1,
		"title":"Test Book",
		"author":"Author",
		"genre":"Fiction",
		"description":"Description",
		"price":19.99
	}`
	assert.JSONEq(t, expected, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestBookHandler_GetAllBooksHandler_Success(t *testing.T) {
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	// Настройка мока
	briefs := []service.BookBrief{
		{
			ID:     1,
			Title:  "Book 1",
			Author: "Author 1",
			Genre:  "Fiction",
			Price:  19.99,
		},
		{
			ID:     2,
			Title:  "Book 2",
			Author: "Author 2",
			Genre:  "Non-Fiction",
			Price:  24.99,
		},
	}
	mockService.On("GetAllBooks", "", 1, 10).Return(briefs, int64(2), nil)

	// Создание запроса
	req, _ := http.NewRequest("GET", "/books", nil)

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.GetAllBooksHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{
		"data": [
			{"id":1, "title":"Book 1", "author":"Author 1", "genre":"Fiction", "price":19.99},
			{"id":2, "title":"Book 2", "author":"Author 2", "genre":"Non-Fiction", "price":24.99}
		],
		"meta": {
			"total":2,
			"page":1,
			"limit":10,
			"totalPages":1
		}
	}`
	assert.JSONEq(t, expected, rr.Body.String())
	mockService.AssertExpectations(t)
}
