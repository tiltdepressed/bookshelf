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

func NewFavouriteService(favService service.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{favService: favService}
}

func (h *FavouriteHandler) AddFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(chi.URLParam(r, "bookID"))
	claims := r.Context().Value("user").(*utils.Claims)
	userID, _ := strconv.Atoi(claims.UserID)

	if err := h.favService.AddFavourite(uint(userID), uint(bookID)); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to add favourite book",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FavouriteHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	bookID, _ := strconv.Atoi(chi.URLParam(r, "bookID"))
	claims := r.Context().Value("user").(*utils.Claims)
	userID, _ := strconv.Atoi(claims.UserID)

	if err := h.favService.RemoveFavourite(uint(userID), uint(bookID)); err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove favourite",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FavouriteHandler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("user").(*utils.Claims)
	userID, _ := strconv.Atoi(claims.UserID)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 15
	}

	books, total, err := h.favService.GetFavourites(uint(userID), page, limit)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get favourites",
		})
		return
	}

	response := map[string]interface{}{
		"data": books,
		"meta": map[string]any{
			"total":      total,
			"page":       page,
			"limit":      limit,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.JSONResponse(w, http.StatusOK, response)
}
