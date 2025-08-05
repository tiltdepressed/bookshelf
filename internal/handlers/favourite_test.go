package handlers_test

import (
	"bookshelf/internal/handlers"
	"bookshelf/internal/service"
	"bookshelf/pkg/utils"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFavouriteService struct {
	mock.Mock
}

func (m *MockFavouriteService) AddFavourite(userID, bookID uint) error {
	args := m.Called(userID, bookID)
	return args.Error(0)
}

func (m *MockFavouriteService) RemoveFavourite(userID, bookID uint) error {
	args := m.Called(userID, bookID)
	return args.Error(0)
}

func (m *MockFavouriteService) GetFavourites(userID uint, page, limit int) ([]service.BookBrief, int64, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]service.BookBrief), args.Get(1).(int64), args.Error(2)
}

func TestFavouriteHandler_AddFavouriteHandler_Success(t *testing.T) {
	mockService := new(MockFavouriteService)
	handler := handlers.NewFavouriteHandler(mockService)

	// Настройка мока
	mockService.On("AddFavourite", uint(1), uint(2)).Return(nil)

	// Создание запроса
	req, _ := http.NewRequest("POST", "/favourites/2", nil)

	// Добавление параметра в роут
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("bookID", "2")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Добавление claims в контекст
	claims := &utils.Claims{
		UserID: "1",
	}
	ctx := context.WithValue(req.Context(), "user", claims)
	req = req.WithContext(ctx)

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.AddFavouriteHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Empty(t, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestFavouriteHandler_GetFavourites_Success(t *testing.T) {
	mockService := new(MockFavouriteService)
	handler := handlers.NewFavouriteHandler(mockService)

	// Настройка мока
	briefs := []service.BookBrief{
		{
			ID:     1,
			Title:  "Fav Book 1",
			Author: "Author 1",
			Genre:  "Fiction",
			Price:  19.99,
		},
	}
	mockService.On("GetFavourites", uint(1), 1, 10).Return(briefs, int64(1), nil)

	// Создание запроса
	req, _ := http.NewRequest("GET", "/favourites?page=1&limit=10", nil)

	// Добавление claims в контекст
	claims := &utils.Claims{
		UserID: "1",
	}
	ctx := context.WithValue(req.Context(), "user", claims)
	req = req.WithContext(ctx)

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.GetFavourites(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{
		"data": [
			{"id":1, "title":"Fav Book 1", "author":"Author 1", "genre":"Fiction", "price":19.99}
		],
		"meta": {
			"total":1,
			"page":1,
			"limit":10,
			"totalPages":1
		}
	}`
	assert.JSONEq(t, expected, rr.Body.String())
	mockService.AssertExpectations(t)
}
