package handlers

import (
	"bookshelf/internal/models"
	"bookshelf/pkg/utils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) RegisterUser(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *MockAuthService) LoginUser(username, password string) (models.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthService) GetUser(currentUserID, targetUserID string) (models.User, error) {
	args := m.Called(currentUserID, targetUserID)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockAuthService) UpdateUserRole(targetUserID, newRole string) (models.User, error) {
	args := m.Called(targetUserID, newRole)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthService) DeleteUser(targetUserID string) error {
	args := m.Called(targetUserID)
	return args.Error(0)
}

func TestAuthHandler_RegisterHandler_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	// Настройка мока
	mockService.On("RegisterUser", "testuser", "password123").Return(nil)

	// Создание запроса
	body := RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.RegisterHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, rr.Code)
	expected := `{"message":"User registered successfully"}`
	assert.JSONEq(t, expected, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestAuthHandler_LoginHandler_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	// Настройка мока
	user := models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Role:     "user",
	}
	mockService.On("LoginUser", "testuser", "password123").Return(user, nil)

	// Создание запроса
	body := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.LoginHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)

	var response LoginResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	mockService.AssertExpectations(t)
}

func TestAuthHandler_GetProfileHandler_Success(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	// Настройка мока
	user := models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Role:     "user",
	}
	mockService.On("GetUser", "1", "1").Return(user, nil)

	// Создание запроса
	req, _ := http.NewRequest("GET", "/users/me", nil)

	// Добавление claims в контекст
	claims := &utils.Claims{
		UserID: "1",
	}
	ctx := context.WithValue(req.Context(), "user", claims)
	req = req.WithContext(ctx)

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.GetProfileHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusOK, rr.Code)
	expected := `{"id":1,"username":"testuser","role":"user"}`
	assert.JSONEq(t, expected, rr.Body.String())
	mockService.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRoleHandler_Forbidden(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	// Создание запроса
	body := UpdateRoleRequest{NewRole: "admin"}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/users/2/role", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Добавление claims в контекст
	claims := &utils.Claims{
		UserID: "1",
		Role:   "user",
	}
	ctx := context.WithValue(req.Context(), "user", claims)
	req = req.WithContext(ctx)

	// Вызов хендлера
	rr := httptest.NewRecorder()
	handler.UpdateUserRoleHandler(rr, req)

	// Проверки
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), "access denied")
}
