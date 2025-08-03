package handlers

import (
	"bookshelf/internal/service"
	"bookshelf/pkg/utils"
	"encoding/json"
	"errors"
	"net/http"

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

func (h *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var req service.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid Request body",
		})
		return
	}

	if err := checkBook(req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]error{
			"error": err,
		})
		return
	}

	book, err := h.bookService.CreateBook(req)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to create book",
		})
		return
	}
	utils.JSONResponse(w, http.StatusCreated, book)
}

func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Book ID is required",
		})
		return
	}

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Book not found",
		})
		return
	}

	utils.JSONResponse(w, http.StatusOK, book)
}

func (h *BookHandler) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Couldn't find books",
		})
		return
	}
	utils.JSONResponse(w, http.StatusOK, books)
}

func (h *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	var req service.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	if err := checkBook(req); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]error{
			"error": err,
		})
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Book ID is required",
		})
		return
	}

	book, err := h.bookService.UpdateBook(id, req)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to update book",
		})
		return
	}

	utils.JSONResponse(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Book ID is required",
		})
		return
	}

	err := h.bookService.DeleteBook(id)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete book",
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
