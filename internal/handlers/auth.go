// Package handlers
package handlers

import (
	"bookshelf/internal/service"
	"bookshelf/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type errorMapping struct {
	errorPattern string
	status       int
}

var errorMappings = []errorMapping{
	{"access denied", http.StatusForbidden},
	{"invalid", http.StatusBadRequest},
	{"required", http.StatusBadRequest},
	{"not found", http.StatusNotFound},
	{"already exists", http.StatusConflict},
	{"credentials", http.StatusUnauthorized},
}

func (h *AuthHandler) handleServiceError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := err.Error()

	for _, mapping := range errorMappings {
		if strings.Contains(strings.ToLower(msg), mapping.errorPattern) {
			status = mapping.status
			break
		}
	}

	utils.JSONResponse(w, status, map[string]string{"error": msg})
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	err := h.authService.RegisterUser(input.Username, input.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	user, err := h.authService.LoginUser(input.Username, input.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	token, err := utils.GenerateToken(
		fmt.Sprintf("%d", user.ID),
		user.Username,
		user.Role,
	)
	if err != nil {
		utils.JSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "User information not found in context",
		})
		return
	}

	user, err := h.authService.GetUser(claims.UserID, claims.UserID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, user)
}

func (h *AuthHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "User information not found in context",
		})
		return
	}

	targetUserID := chi.URLParam(r, "id")
	if _, err := strconv.ParseUint(targetUserID, 10, 64); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format",
		})
		return
	}
	if targetUserID == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
		return
	}

	user, err := h.authService.GetUser(claims.UserID, targetUserID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}

func (h *AuthHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "User information not found in context",
		})
		return
	}

	users, err := h.authService.GetAllUsers(claims.UserID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, users)
}

func (h *AuthHandler) UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "User information not found in context",
		})
		return
	}

	targetUserID := chi.URLParam(r, "id")
	if targetUserID == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "User ID is required",
		})
		return
	}

	var input struct {
		NewRole string `json:"new_role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	updatedUser, err := h.authService.UpdateUserRole(claims.UserID, targetUserID, input.NewRole)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, updatedUser)
}

func (h *AuthHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claims)
	if !ok {
		utils.JSONResponse(w, http.StatusUnauthorized, map[string]string{
			"error": "User information not found in context",
		})
		return
	}

	targetUserID := chi.URLParam(r, "id")
	if _, err := strconv.ParseUint(targetUserID, 10, 64); err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format",
		})
		return
	}
	if targetUserID == "" {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "User ID os required",
		})
		return
	}

	err := h.authService.DeleteUser(claims.UserID, targetUserID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
