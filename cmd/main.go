package main

import (
	"bookshelf/internal/config/db"
	"bookshelf/internal/handlers"
	"bookshelf/internal/middleware" // Локальный middleware
	"bookshelf/internal/repository"
	"bookshelf/internal/service"
	"bookshelf/pkg/utils"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware" // Псевдоним для избежания конфликта
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to db: %s", err.Error())
	}

	utils.InitJWT()

	// Инициализация слоёв
	authRepo := repository.NewAuthRepository(database)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	bookRepo := repository.NewBookRepository(database)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Публичные роуты
	r.Group(func(r chi.Router) {
		r.Post("/api/auth/register", authHandler.RegisterHandler)
		r.Post("/api/auth/login", authHandler.LoginHandler)

		r.Get("/api/books", bookHandler.GetAllBooksHandler)
		r.Get("/api/books/{id}", bookHandler.GetBookByIDHandler)
	})

	// Защищенные роуты (для всех авторизованных)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)

		r.Get("/api/users/me", authHandler.GetProfileHandler) // Мой профиль
		r.Get("/api/users/{id}", authHandler.GetUserHandler)  // Профиль другого пользователя
	})

	// Админские роуты (только для админов)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware, middleware.AdminOnlyMiddleware)

		r.Get("/api/users", authHandler.GetAllUsersHandler)                // Все пользователи
		r.Patch("/api/users/{id}/role", authHandler.UpdateUserRoleHandler) // Изменение роли
		r.Delete("/api/users/{id}", authHandler.DeleteUserHandler)         // Удаление пользователя

		r.Post("/api/books", bookHandler.CreateBookHandler)
		r.Patch("/api/books/{id}", bookHandler.UpdateBookHandler)
		r.Delete("/api/books/{id}", bookHandler.DeleteBookHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started and listening on port %s", port)
	log.Printf("Try: http://localhost:%s", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Could not start listening: %s", err.Error())
	}
}
