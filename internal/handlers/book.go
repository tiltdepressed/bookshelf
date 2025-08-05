package handlers

import (
	"bookshelf/internal/service"
	"bookshelf/pkg/utils"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) *BookHandler {
	return &BookHandler{bookService: bookService}
}

func checkBook(book service.BookRequest) error {
	if book.Author == "" || book.Title == "" || book.Genre == "" || book.Description == "" || book.Price < 0 {
		return errors.New("invalid book data")
	}
	return nil
}

// CreateBookHandler godoc
// @Summary Создание новой книги
// @Description Создание новой книги в системе
// @Tags Books
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param input body service.BookRequest true "Данные книги"
// @Success 201 {object} BookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var req service.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid request body"})
		return
	}

	if err := checkBook(req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	book, err := h.bookService.CreateBook(req)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Failed to create book"})
		return
	}

	utils.JSONResponse(w, http.StatusCreated, BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Description: book.Description,
		Price:       book.Price,
	})
}

// GetBookByIDHandler godoc
// @Summary Получение книги по ID
// @Description Получение информации о книге по её идентификатору
// @Tags Books
// @Produce json
// @Param id path string true "ID книги"
// @Success 200 {object} BookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Book ID is required"})
		return
	}

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusNotFound, ErrorResponse{"Book not found"})
		return
	}

	utils.JSONResponse(w, http.StatusOK, BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Description: book.Description,
		Price:       book.Price,
	})
}

// GetAllBooksHandler godoc
// @Summary Получение списка книг
// @Description Получение списка книг с возможностью фильтрации по жанру и пагинацией
// @Tags Books
// @Produce json
// @Param genre query string false "Фильтр по жанру"
// @Param page query int false "Номер страницы (по умолчанию 1)" default(1)
// @Param limit query int false "Количество книг на странице (по умолчанию 10, максимум 100)" default(10)
// @Success 200 {object} PaginatedBooksResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	books, total, err := h.bookService.GetAllBooks(genre, page, limit)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Couldn't find books"})
		return
	}

	var bookResponses []BookBriefResponse
	for _, book := range books {
		bookResponses = append(bookResponses, BookBriefResponse{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Genre:  book.Genre,
			Price:  book.Price,
		})
	}

	response := PaginatedBooksResponse{
		Data: bookResponses,
		Meta: struct {
			Total      int64 `json:"total" example:"100"`
			Page       int   `json:"page" example:"1"`
			Limit      int   `json:"limit" example:"10"`
			TotalPages int   `json:"totalPages" example:"10"`
		}{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		},
	}

	utils.JSONResponse(w, http.StatusOK, response)
}

// GetAllGenresHandler godoc
// @Summary Получение списка жанров
// @Description Получение списка всех доступных жанров книг
// @Tags Books
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} ErrorResponse
// @Router /books/genres [get]
func (h *BookHandler) GetAllGenresHandler(w http.ResponseWriter, r *http.Request) {
	genres, err := h.bookService.GetAllGenres()
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Couldn't get genres"})
		return
	}

	utils.JSONResponse(w, http.StatusOK, genres)
}

// UpdateBookHandler godoc
// @Summary Обновление информации о книге
// @Description Обновление информации о существующей книге
// @Tags Books
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID книги"
// @Param input body service.BookRequest true "Обновленные данные книги"
// @Success 200 {object} BookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	var req service.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid request body"})
		return
	}

	if err := checkBook(req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{err.Error()})
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Book ID is required"})
		return
	}

	book, err := h.bookService.UpdateBook(id, req)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Failed to update book"})
		return
	}

	utils.JSONResponse(w, http.StatusOK, BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Genre:       book.Genre,
		Description: book.Description,
		Price:       book.Price,
	})
}

// DeleteBookHandler godoc
// @Summary Удаление книги
// @Description Удаление книги из системы
// @Tags Books
// @Security ApiKeyAuth
// @Param id path string true "ID книги"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Book ID is required"})
		return
	}

	err := h.bookService.DeleteBook(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Failed to delete book"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
