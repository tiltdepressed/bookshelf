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

// RegisterHandler godoc
// @Summary Регистрация нового пользователя
// @Description Создание нового аккаунта пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
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

// LoginHandler godoc
// @Summary Аутентификация пользователя
// @Description Вход пользователя в систему и получение токена
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
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

// GetProfileHandler godoc
// @Summary Получение профиля текущего пользователя
// @Description Получение информации о текущем аутентифицированном пользователе
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/me [get]
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

// GetUserHandler godoc
// @Summary Получение информации о пользователе
// @Description Получение информации о пользователе по ID (доступно администраторам)
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
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

// GetAllUsersHandler godoc
// @Summary Получение списка всех пользователей
// @Description Получение списка всех пользователей (доступно администраторам)
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /users [get]
func (h *AuthHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, users)
}

// UpdateUserRoleHandler godoc
// @Summary Изменение роли пользователя
// @Description Изменение роли пользователя (доступно администраторам)
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param input body UpdateRoleRequest true "Новая роль"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id}/role [put]
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

	if claims.UserID == targetUserID {
		utils.JSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "Cannot change your own role",
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

	updatedUser, err := h.authService.UpdateUserRole(targetUserID, input.NewRole)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, updatedUser)
}

// DeleteUserHandler godoc
// @Summary Удаление пользователя
// @Description Удаление пользователя (доступно администраторам)
// @Tags Users
// @Security ApiKeyAuth
// @Param id path string true "ID пользователя"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (h *AuthHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	err := h.authService.DeleteUser(targetUserID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
