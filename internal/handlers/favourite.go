package handlers

import (
	"bookshelf/internal/service"
	"bookshelf/pkg/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type FavouriteHandler struct {
	favService service.FavouriteService
}

func NewFavouriteHandler(favService service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{favService: favService}
}

// AddFavouriteHandler godoc
// @Summary Добавление книги в избранное
// @Description Добавляет книгу в список избранных для текущего пользователя
// @Tags Favourites
// @Security ApiKeyAuth
// @Param bookID path int true "ID книги"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /favourites/{bookID} [post]
func (h *FavouriteHandler) AddFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid book ID"})
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, ErrorResponse{"User information not found"})
		return
	}

	userID, _ := strconv.Atoi(claims.UserID)

	if err := h.favService.AddFavourite(uint(userID), uint(bookID)); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Failed to add favourite book"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveFavourite godoc
// @Summary Удаление книги из избранного
// @Description Удаляет книгу из списка избранных для текущего пользователя
// @Tags Favourites
// @Security ApiKeyAuth
// @Param bookID path int true "ID книги"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /favourites/{bookID} [delete]
func (h *FavouriteHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	bookID, err := strconv.Atoi(chi.URLParam(r, "bookID"))
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, ErrorResponse{"Invalid book ID"})
		return
	}

	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, ErrorResponse{"User information not found"})
		return
	}

	userID, _ := strconv.Atoi(claims.UserID)

	if err := h.favService.RemoveFavourite(uint(userID), uint(bookID)); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Failed to remove favourite"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetFavourites godoc
// @Summary Получение списка избранных книг
// @Description Возвращает список избранных книг для текущего пользователя с пагинацией
// @Tags Favourites
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)" default(1)
// @Param limit query int false "Количество книг на странице (по умолчанию 10, максимум 100)" default(10)
// @Success 200 {object} PaginatedBooksResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /favourites [get]
func (h *FavouriteHandler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, ErrorResponse{"User information not found"})
		return
	}

	userID, _ := strconv.Atoi(claims.UserID)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	books, total, err := h.favService.GetFavourites(uint(userID), page, limit)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, ErrorResponse{"Failed to get favourites"})
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
			TotalPages: int((total + int64(limit) - 1) / int64(limit)),
		},
	}

	utils.JSONResponse(w, http.StatusOK, response)
}
