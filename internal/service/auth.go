// Package service
package service

import (
	"bookshelf/internal/models"
	"bookshelf/internal/repository"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userContext struct {
	user models.User
	err  error
}

type AuthService interface {
	RegisterUser(username, password string) error
	LoginUser(username, password string) (models.User, error)
	// Полудаминский метод: простые смертные могут смотреть только свой профиль, админ - любые
	GetUser(currentUserID string, targetUserID string) (models.User, error)
	// Админские методы
	GetAllUsers(currentUserID string) ([]models.User, error)
	UpdateUserRole(currentUserID, targetUserID, newRole string) (models.User, error)
	DeleteUser(currentUserID string, targetUserID string) error
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) getUserContext(userID string) *userContext {
	user, err := s.repo.GetUserByID(userID)
	return &userContext{user, err}
}

func (s *authService) RegisterUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password required")
	}
	_, err := s.repo.GetUserByUsername(username)
	if errors.Is(err, nil) {
		return errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	return s.repo.CreateUser(newUser)
}

func (s *authService) LoginUser(username, password string) (models.User, error) {
	if username == "" || password == "" {
		return models.User{}, errors.New("username and password required")
	}

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("invalid credentials")
		}
		return models.User{}, fmt.Errorf("database error: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		time.Sleep(2 * time.Second) // Замедление атак перебора
		return models.User{}, errors.New("invalid credentials")
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *authService) GetUser(currentUserID string, targetUserID string) (models.User, error) {
	ctx := s.getUserContext(currentUserID)
	if ctx.err != nil {
		return models.User{}, ctx.err
	}

	if fmt.Sprintf("%d", ctx.user.ID) != targetUserID && ctx.user.Role != "admin" {
		return models.User{}, errors.New("access denied")
	}

	user, err := s.repo.GetUserByID(targetUserID)
	if err != nil {
		return models.User{}, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *authService) GetAllUsers(currentUserID string) ([]models.User, error) {
	ctx := s.getUserContext(currentUserID)
	if ctx.err != nil {
		return []models.User{}, ctx.err
	}

	if ctx.user.Role != "admin" {
		return []models.User{}, errors.New("access denied: admin privileges required")
	}

	users, err := s.repo.GetAllUsers()
	if err != nil {
		return []models.User{}, fmt.Errorf("failed to get users: %w", err)
	}

	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

func (s *authService) UpdateUserRole(currentUserID, targetUserID, newRole string) (models.User, error) {
	if newRole != "admin" && newRole != "user" {
		return models.User{}, errors.New("invalid role, must be 'admin' or 'user'")
	}

	ctx := s.getUserContext(currentUserID)
	if ctx.err != nil {
		return models.User{}, ctx.err
	}

	if ctx.user.Role != "admin" {
		return models.User{}, errors.New("access denied: admin privileges required")
	}

	user, err := s.repo.GetUserByID(targetUserID)
	if err != nil {
		return models.User{}, err
	}

	if currentUserID == targetUserID {
		return models.User{}, errors.New("cannot change own role")
	}

	user.Role = newRole
	if err := s.repo.UpdateUser(user); err != nil {
		return models.User{}, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *authService) DeleteUser(currentUserID string, targetUserID string) error {
	ctx := s.getUserContext(currentUserID)
	if ctx.err != nil {
		return ctx.err
	}

	if ctx.user.Role != "admin" {
		return errors.New("access denied: admin privileges required")
	}

	return s.repo.DeleteUser(targetUserID)
}
