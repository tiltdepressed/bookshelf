package middleware

import (
	"bookshelf/pkg/utils"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	// Инициализация JWT
	utils.InitJWT()
	token, _ := utils.GenerateToken("1", "testuser", "user")

	// Создание запроса с токеном
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Создание тестового обработчика
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user").(*utils.Claims)
		assert.True(t, ok)
		assert.Equal(t, "1", claims.UserID)
		w.WriteHeader(http.StatusOK)
	})

	// Применение middleware
	handler := JWTAuthMiddleware(nextHandler)

	// Вызов
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	// Создание запроса с невалидным токеном
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	// Создание тестового обработчика
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	// Применение middleware
	handler := JWTAuthMiddleware(nextHandler)

	// Вызов
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid token")
}

func TestAdminOnlyMiddleware_AdminUser(t *testing.T) {
	// Создание запроса
	req, _ := http.NewRequest("GET", "/admin", nil)

	// Добавление claims в контекст
	claims := &utils.Claims{
		Role: "admin",
	}
	ctx := context.WithValue(req.Context(), "user", claims)
	req = req.WithContext(ctx)

	// Создание тестового обработчика
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Применение middleware
	handler := AdminOnlyMiddleware(nextHandler)

	// Вызов
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAdminOnlyMiddleware_NonAdminUser(t *testing.T) {
	// Создание запроса
	req, _ := http.NewRequest("GET", "/admin", nil)

	// Добавление claims в контекст
	claims := &utils.Claims{
		Role: "user",
	}
	ctx := context.WithValue(req.Context(), "user", claims)
	req = req.WithContext(ctx)

	// Создание тестового обработчика
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Next handler should not be called")
	})

	// Применение middleware
	handler := AdminOnlyMiddleware(nextHandler)

	// Вызов
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Проверки
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), "Admin privileges required")
}
