package main

import (
	"bookshelf/internal/config/db"
	"bookshelf/internal/handlers"
	"bookshelf/internal/middleware"
	"bookshelf/internal/repository"
	"bookshelf/internal/service"
	"bookshelf/pkg/cache"
	"bookshelf/pkg/utils"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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

	redisCache, err := cache.NewRedisCache(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %s", err.Error())
	}

	utils.InitJWT()

	// Инициализация слоёв
	authRepo := repository.NewAuthRepository(database)
	authService := service.NewAuthService(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	bookRepo := repository.NewBookRepository(database)
	bookService := service.NewBookService(bookRepo, redisCache)
	bookHandler := handlers.NewBookHandler(bookService)

	favRepo := repository.NewFavouriteRepository(database)
	favService := service.NewFavouriteService(favRepo)
	favHandler := handlers.NewFavouriteHandler(favService)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Публичные роуты
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", authHandler.RegisterHandler)
		r.Post("/auth/login", authHandler.LoginHandler)

		r.Get("/books", bookHandler.GetAllBooksHandler)
		r.Get("/books/{id}", bookHandler.GetBookByIDHandler)

		r.Get("/books/genres", bookHandler.GetAllGenresHandler)
	})

	// Защищенные роуты (для всех авторизованных)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware)

		r.Get("/users/me", authHandler.GetProfileHandler) // Мой профиль
		r.Get("/users/{id}", authHandler.GetUserHandler)  // Профиль другого пользователя

		r.Get("/favourites/me", favHandler.GetFavourites)
		r.Post("/favourites/{bookID}", favHandler.AddFavouriteHandler)
		r.Delete("/favourites/{bookID}", favHandler.RemoveFavourite)
	})

	// Админские роуты (только для админов)
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware, middleware.AdminOnlyMiddleware)

		r.Get("/users", authHandler.GetAllUsersHandler)
		r.Patch("/users/{id}/role", authHandler.UpdateUserRoleHandler)
		r.Delete("/users/{id}", authHandler.DeleteUserHandler)

		r.Post("/books", bookHandler.CreateBookHandler)
		r.Patch("/books/{id}", bookHandler.UpdateBookHandler)
		r.Delete("/books/{id}", bookHandler.DeleteBookHandler)
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
